// +build integration

package keys

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/DSuhinin/passbase-test-task/core"
	"github.com/DSuhinin/passbase-test-task/core/errors"

	"github.com/DSuhinin/passbase-test-task/app"
	"github.com/DSuhinin/passbase-test-task/app/api/response"
	"github.com/DSuhinin/passbase-test-task/app/config"
	"github.com/DSuhinin/passbase-test-task/test/fixtures"
	"github.com/DSuhinin/passbase-test-task/test/helper"
)

type RegenerateKeyTestSuite struct {
	suite.Suite
	fixtures *fixtures.Fixtures
}

// TestRegenerateKey is an entry point to run all tests in current Test Suite.
func TestRegenerateKey(t *testing.T) {
	suite.Run(t, new(RegenerateKeyTestSuite))
}

// SetupSuite prepare everything for tests.
func (s *RegenerateKeyTestSuite) SetupSuite() {
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

// TestRegenerateKey_OK makes test of `GET /keys/:key_id/regenerate` for success case.
func (s *RegenerateKeyTestSuite) TestRegenerateKey_OK() {
	assert.Nil(s.T(), s.fixtures.LoadRegenerateKeyData())
	defer func() {
		assert.Nil(s.T(), s.fixtures.UnloadFixtures())
	}()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPut,
		helper.StrReplace(
			fmt.Sprintf("%s%s", os.Getenv("SERVICE_BASE_URL"), app.RegenerateKeyRoute),
			[]string{":key_id"},
			[]interface{}{1},
		),
		nil,
	)
	req.Header.Set("Authorization", fmt.Sprintf("AdminKey %s", "supersecurekey"))
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

	response := response.Key{}
	assert.Nil(s.T(), json.Unmarshal(body, &response))
	assert.Equal(s.T(), 1, response.ID)
	assert.NotEmpty(s.T(), response.Value)

	_, err = uuid.FromString(response.Value)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "e12a1983-046a-4f2c-b5a2-e27a6851ec4c", response.Value)
}

// TestRegenerateKey_AuthError makes test of `GET /keys/:key_id/regenerate` for Auth error case.
func (s *CreateKeyTestSuite) TestRegenerateKey_AuthError() {
	assert.Nil(s.T(), s.fixtures.LoadRegenerateKeyData())
	defer func() {
		assert.Nil(s.T(), s.fixtures.UnloadFixtures())
	}()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPut,
		helper.StrReplace(
			fmt.Sprintf("%s%s", os.Getenv("SERVICE_BASE_URL"), app.RegenerateKeyRoute),
			[]string{":key_id"},
			[]interface{}{1},
		),
		nil,
	)
	req.Header.Set("Authorization", "AdminKey incorrectkey")
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
