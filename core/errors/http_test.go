// +build unit

package errors

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

const testRoute = "/test"

func CustomError1() error { return NewHTTPInternal(10000, "something went wrong") }

var CustomError2 = NewHTTPInternal(10000, "something went wrong")

// dummy function to test logging the stacktrace
func doWork1() error {
	return CustomError1()
}

// dummy function to test logging the stacktrace
func doWork2() error {
	return CustomError2
}

// dummy function to test logging the stacktrace
func doWork3() error {
	return CustomError2.WithError(New("wrapped error happened"))
}

func TestSetHTTPError_CustomError1_OK(t *testing.T) {
	// Given
	testHandler := func(c *gin.Context) {
		err := doWork1()
		SetHTTPError(c, err)
		return
	}
	testRouter := setupTestRouter(testHandler)

	logHook := logtest.NewGlobal()

	// When
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, testRoute, nil)
	testRouter.ServeHTTP(w, req)

	// Then
	// HTTP status code is correct
	require.Equal(t, http.StatusInternalServerError, w.Code)
	// HTTP body contains the custom code and custom message
	require.JSONEq(t, `{"code": 10000, "message": "something went wrong"}`, w.Body.String())

	// One error log line was logged
	require.Equal(t, 1, len(logHook.Entries))
	require.Equal(t, logrus.ErrorLevel, logHook.LastEntry().Level)
	// log line contains the stack trace fragment
	message, err := logHook.LastEntry().String()
	require.Nil(t, err)
	require.True(t, strings.Contains(message, "doWork1"))
	// log line contains the error message
	require.True(t, strings.Contains(logHook.LastEntry().Message, "something went wrong"))
}

func TestSetHTTPError_CustomError2_NoStacktrace(t *testing.T) {
	// Given
	testHandler := func(c *gin.Context) {
		err := doWork2()
		SetHTTPError(c, err)
		return
	}
	testRouter := setupTestRouter(testHandler)

	logHook := logtest.NewGlobal()

	// When
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, testRoute, nil)
	testRouter.ServeHTTP(w, req)

	// Then
	// HTTP status code is correct
	require.Equal(t, http.StatusInternalServerError, w.Code)
	// HTTP body contains the custom code and custom message
	require.JSONEq(t, `{"code": 10000, "message": "something went wrong"}`, w.Body.String())

	// One error log line was logged
	require.Equal(t, 1, len(logHook.Entries))
	require.Equal(t, logrus.ErrorLevel, logHook.LastEntry().Level)
	// log line does not contain the stack trace fragment
	require.False(t, strings.Contains(logHook.LastEntry().Message, "doWork2"))
	// log line contains the error message
	require.True(t, strings.Contains(logHook.LastEntry().Message, "something went wrong"))
}

func TestSetHTTPError_CustomError2_WithError_OK(t *testing.T) {
	// Given
	testHandler := func(c *gin.Context) {
		err := doWork3()
		SetHTTPError(c, err)
		return
	}
	testRouter := setupTestRouter(testHandler)

	logHook := logtest.NewGlobal()

	// When
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, testRoute, nil)
	testRouter.ServeHTTP(w, req)

	// Then
	// HTTP status code is correct
	require.Equal(t, http.StatusInternalServerError, w.Code)
	// HTTP body contains the custom code and custom message
	require.JSONEq(t, `{"code": 10000, "message": "something went wrong"}`, w.Body.String())

	// One error log line was logged
	require.Equal(t, 1, len(logHook.Entries))
	require.Equal(t, logrus.ErrorLevel, logHook.LastEntry().Level)
	// log line does not contain the stack trace fragment
	message, err := logHook.LastEntry().String()
	require.Nil(t, err)
	require.True(t, strings.Contains(message, "doWork3"))
	// log line contains the error message
	require.True(t, strings.Contains(logHook.LastEntry().Message, "wrapped error happened"))
}

func setupTestRouter(handler func(c *gin.Context)) *gin.Engine {
	r := gin.Default()
	r.GET(testRoute, handler)
	return r
}
