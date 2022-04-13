package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// SetHTTPError sets the HTTP error in Gin context.
// A return statement must be made in the Gin handler function afterwards.
func SetHTTPError(ctx *gin.Context, err error) {
	httpError, ok := err.(*HTTP)
	if !ok {
		httpError = NewHTTPInternal(http.StatusInternalServerError, "unmapped internal server error")
		httpError.OriginalError = err
	}

	if httpError.OriginalError != nil {
		logEntry := log.WithField(
			"stack", strings.Split(fmt.Sprintf("%+s", httpError.OriginalError), "\n\t"),
		)
		switch dbError := Cause(httpError.OriginalError).(type) {
		case *pq.Error:
			logEntry.Errorf("code: %s, error: %+v, details: %s", dbError.Code, dbError.Message, dbError.Detail)
		default:
			logEntry.Errorf("%s", httpError.OriginalError) // logs stacktrace too, if error was created using this package
		}
	}

	ctx.JSON(httpError.HTTPCode, httpError)
}

// HTTP Generic error response
type HTTP struct {
	HTTPCode      int    `json:"-"`
	Code          int    `json:"code,omitempty"`
	Message       string `json:"message,omitempty"`
	OriginalError error  `json:"-"`
}

// Error implements Go error interface.
// Displays code and message as JSON.
func (e *HTTP) Error() string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`, e.Code, e.Message)
}

// WithError make possible to store original error object.
func (e *HTTP) WithError(err error) *HTTP {
	e.OriginalError = err

	return e
}

// NewHTTPError creates new HTTP error with the given HTTP status code.
// If used in combination with SetHTTPError, the given HTTP status will be returned and parameters code and message
// will be displayed in the HTTP response body.
// Stack trace is wrapped into the HTTP error, which is unwrapped and logged when passed to SetHTTPError.
// Wrapped error can be overwritten using NewHTTPError().WithError(err).
func NewHTTPError(httpCode int, code int, message string) *HTTP {
	return &HTTP{
		HTTPCode:      httpCode,
		Code:          code,
		Message:       message,
		OriginalError: New(message),
	}
}

// NewHTTPBadRequest creates new HTTP error with HTTP status 400.
// If used in combination with SetHTTPError, the given HTTP status will be returned and parameters code and message
// will be displayed in the HTTP response body.
// Stack trace is wrapped into the HTTP error, which is unwrapped and logged when passed to SetHTTPError.
// Wrapped error can be overwritten using NewHTTPError().WithError(err).
func NewHTTPBadRequest(code int, message string) *HTTP {
	return NewHTTPError(http.StatusBadRequest, code, message)
}

// NewHTTPUnauthorized creates new HTTP error with HTTP status 401.
// If used in combination with SetHTTPError, the given HTTP status will be returned and parameters code and message
// will be displayed in the HTTP response body.
// Stack trace is wrapped into the HTTP error, which is unwrapped and logged when passed to SetHTTPError.
// Wrapped error can be overwritten using NewHTTPError().WithError(err).
func NewHTTPUnauthorized(code int, message string) *HTTP {
	return NewHTTPError(http.StatusUnauthorized, code, message)
}

// NewHTTPNotFound creates new HTTP error with HTTP status 404.
// If used in combination with SetHTTPError, the given HTTP status will be returned and parameters code and message
// will be displayed in the HTTP response body.
// Stack trace is wrapped into the HTTP error, which is unwrapped and logged when passed to SetHTTPError.
// Wrapped error can be overwritten using NewHTTPError().WithError(err).
func NewHTTPNotFound(code int, message string) *HTTP {
	return NewHTTPError(http.StatusNotFound, code, message)
}

// NewHTTPInternal creates new HTTP error with HTTP status 500.
// If used in combination with SetHTTPError, the given HTTP status will be returned and parameters code and message
// will be displayed in the HTTP response body.
// Stack trace is wrapped into the HTTP error, which is unwrapped and logged when passed to SetHTTPError.
// Wrapped error can be overwritten using NewHTTPError().WithError(err).
func NewHTTPInternal(code int, message string) *HTTP {
	return NewHTTPError(http.StatusInternalServerError, code, message)
}
