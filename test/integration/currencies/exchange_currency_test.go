// +build integration

package currencies

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/DSuhinin/passbase-test-task/core"
	"github.com/DSuhinin/passbase-test-task/core/errors"

	"github.com/DSuhinin/passbase-test-task/app"
	"github.com/DSuhinin/passbase-test-task/app/api/response"
	"github.com/DSuhinin/passbase-test-task/app/config"
	"github.com/DSuhinin/passbase-test-task/test/data"
	"github.com/DSuhinin/passbase-test-task/test/fixtures"
)

type ExchangeCurrencyTestSuite struct {
	suite.Suite
	fixtures *fixtures.Fixtures
}

// TestExchangeCurrency is an entry point to run all tests in current Test Suite.
func TestExchangeCurrency(t *testing.T) {
	suite.Run(t, new(ExchangeCurrencyTestSuite))
}

// SetupSuite prepare everything for tests.
func (s *ExchangeCurrencyTestSuite) SetupSuite() {
	// 1. init config.
	appConfig, err := config.New()
	assert.Nil(s.T(), err)

	// 2.initialize db connections.
	dbConnection, err := core.NewDB().GetConnection(
		appConfig.DatabaseUser,
		appConfig.DatabasePass,
		core.PostgresType,
		appConfig.DatabaseName,
		appConfig.DatabaseHost,
	)
	assert.Nil(s.T(), err)

	// 3. init fixtures.
	s.fixtures = fixtures.NewFixtures(dbConnection)
}

// TestExchangeCurrency_OK makes test of `GET /currencies/exchange` for success case.
func (s *ExchangeCurrencyTestSuite) TestExchangeCurrency_OK() {
	assert.Nil(s.T(), s.fixtures.LoadCurrencyExchangeData())
	defer func() {
		assert.Nil(s.T(), s.fixtures.UnloadFixtures())
	}()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf("%s%s?from=USD&to=EUR&amount=100", os.Getenv("SERVICE_BASE_URL"), app.CurrenciesExchangeRoute),
		nil,
	)
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", "e12a1983-046a-4f2c-b5a2-e27a6851ec4c"))
	assert.Nil(s.T(), err)

	client := http.Client{}
	resp, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp.Body)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(s.T(), err)
	// nolint
	defer resp.Body.Close()

	response := response.CurrencyExchange{}
	assert.Nil(s.T(), json.Unmarshal(body, &response))
	assert.Equal(s.T(), data.CurrencyExchangeResponse, response)
}

// TestExchangeCurrency_ValidationError makes test of `GET /currencies/exchange` for different validation error cases.
func (s *ExchangeCurrencyTestSuite) TestExchangeCurrency_ValidationError() {

	assert.Nil(s.T(), s.fixtures.LoadCurrencyExchangeData())
	defer func() {
		assert.Nil(s.T(), s.fixtures.UnloadFixtures())
	}()

	type test struct {
		name         string
		query        string
		errorCode    int
		errorMessage string
	}

	tableTest := []test{
		{
			name:         "WithIncorrectFromParameter",
			query:        "from=incorrect",
			errorCode:    20000,
			errorMessage: "`from` parameter should be USD or EUR",
		},
		{
			name:         "WithIncorrectToParameter",
			query:        "from=USD&to=incorrect",
			errorCode:    20010,
			errorMessage: "`to` parameter should be USD or EUR",
		},
		{
			name:         "WithIncorrectAmountParameter",
			query:        "from=USD&to=EUR&amount=",
			errorCode:    20020,
			errorMessage: "`amount` parameter should be grater then zero",
		},
	}

	//nolint:scopelint
	for _, tt := range tableTest {
		s.T().Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				fmt.Sprintf("%s%s?%s", os.Getenv("SERVICE_BASE_URL"), app.CurrenciesExchangeRoute, tt.query),
				nil,
			)
			req.Header.Set("Authorization", fmt.Sprintf("Key %s", "e12a1983-046a-4f2c-b5a2-e27a6851ec4c"))
			assert.Nil(s.T(), err)

			client := http.Client{}
			resp, err := client.Do(req)
			assert.Nil(s.T(), err)
			assert.NotNil(s.T(), resp.Body)
			assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

			body, err := ioutil.ReadAll(resp.Body)
			assert.Nil(s.T(), err)
			// nolint
			defer resp.Body.Close()

			response := errors.HTTP{}
			assert.Nil(s.T(), json.Unmarshal(body, &response))
			assert.Equal(s.T(), tt.errorCode, response.Code)
			assert.Equal(s.T(), tt.errorMessage, response.Message)
		})
	}
}

// TestExchangeCurrency_AuthError makes test of `GET /currencies/exchange` for Auth error case.
func (s *ExchangeCurrencyTestSuite) TestExchangeCurrency_AuthError() {
	assert.Nil(s.T(), s.fixtures.LoadCurrencyExchangeData())
	defer func() {
		assert.Nil(s.T(), s.fixtures.UnloadFixtures())
	}()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf("%s%s?from=USD&to=EUR&amount=100", os.Getenv("SERVICE_BASE_URL"), app.CurrenciesExchangeRoute),
		nil,
	)
	req.Header.Set("Authorization", "Key incorrectkey")
	assert.Nil(s.T(), err)

	client := http.Client{}
	resp, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp.Body)
	assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(s.T(), err)
	// nolint
	defer resp.Body.Close()

	response := errors.HTTP{}
	assert.Nil(s.T(), json.Unmarshal(body, &response))
	assert.Equal(s.T(), 100001, response.Code)
	assert.Equal(s.T(), "unauthorized request", response.Message)
}
