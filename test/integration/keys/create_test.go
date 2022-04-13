// +build integration

package campaign

import (
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

	"github.com/DSuhinin/passbase-test-task/app"
	"github.com/DSuhinin/passbase-test-task/app/api/response"
	"github.com/DSuhinin/passbase-test-task/app/config"
	"github.com/DSuhinin/passbase-test-task/test/fixtures"
)

type CreateKeyTestSuite struct {
	suite.Suite
	fixtures *fixtures.Fixtures
}

// TestCreateKey is an entry point to run all tests in current Test Suite.
func TestCreateKey(t *testing.T) {
	suite.Run(t, new(CreateKeyTestSuite))
}

// SetupSuite prepare everything for tests.
func (s *CreateKeyTestSuite) SetupSuite() {
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

// TestCreateKey_OK makes test of `POST /keys` for success case.
func (s *CreateKeyTestSuite) TestCreateKey_OK() {
	defer func() {
		assert.Nil(s.T(), s.fixtures.UnloadFixtures())
	}()

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s", os.Getenv("SERVICE_BASE_URL"), app.CreateKeyRoute),
		nil,
	)
	req.Header.Set("Authorization", fmt.Sprintf("AdminKey %s", "supersecurekey"))
	assert.Nil(s.T(), err)

	client := http.Client{}
	resp, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(s.T(), err)
	// nolint
	defer resp.Body.Close()

	response := response.Key{}
	assert.Nil(s.T(), json.Unmarshal(body, &response))
	assert.Equal(s.T(), 1, response.ID)
	assert.NotEmpty(s.T(), response.Key)

	_, err = uuid.FromString(response.Key)
	assert.Nil(s.T(), err)
}
