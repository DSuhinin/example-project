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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/DSuhinin/passbase-test-task/core"
	"github.com/DSuhinin/passbase-test-task/core/errors"

	"github.com/DSuhinin/passbase-test-task/app"
	"github.com/DSuhinin/passbase-test-task/app/config"
	"github.com/DSuhinin/passbase-test-task/test/fixtures"
	"github.com/DSuhinin/passbase-test-task/test/helper"
)

type DeleteKeyTestSuite struct {
	suite.Suite
	fixtures *fixtures.Fixtures
}

// TestDeleteKey is an entry point to run all tests in current Test Suite.
func TestDeleteKey(t *testing.T) {
	suite.Run(t, new(DeleteKeyTestSuite))
}

// SetupSuite prepare everything for tests.
func (s *DeleteKeyTestSuite) SetupSuite() {
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

// TestDeleteKey_OK makes test of `DELETE /keys/:key_id` for success case.
func (s *DeleteKeyTestSuite) TestDeleteKey_OK() {
	assert.Nil(s.T(), s.fixtures.LoadDeleteKeyData())
	defer func() {
		assert.Nil(s.T(), s.fixtures.UnloadFixtures())
	}()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		helper.StrReplace(
			fmt.Sprintf("%s%s", os.Getenv("SERVICE_BASE_URL"), app.DeleteKeyRoute),
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
	//nolint
	defer resp.Body.Close()
	assert.Equal(s.T(), http.StatusNoContent, resp.StatusCode)
}

// TestDeleteKey_AuthError makes test of `DELETE /keys/:key_id` for Auth error case.
func (s *CreateKeyTestSuite) TestDeleteKey_AuthError() {
	assert.Nil(s.T(), s.fixtures.LoadDeleteKeyData())
	defer func() {
		assert.Nil(s.T(), s.fixtures.UnloadFixtures())
	}()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		helper.StrReplace(
			fmt.Sprintf("%s%s", os.Getenv("SERVICE_BASE_URL"), app.DeleteKeyRoute),
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
