package middleware

import (
	"github.com/DSuhinin/passbase-test-task/app/service/keys/dao"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/DSuhinin/passbase-test-task/core/errors"
)

// KeyIDPlaceholder is a Value ID placeholder.
const KeyIDPlaceholder = "key_id"

// ValidateKeyID validates `:key_id` placeholder from the request.
func ValidateKeyID(context *gin.Context) {
	value, err := strconv.Atoi(context.Param(KeyIDPlaceholder))
	if err != nil {
		errors.SetHTTPError(
			context, errors.PathParametersParsingError.WithError(
				errors.Wrapf(err, "wrong format of :%s", KeyIDPlaceholder),
			),
		)
		context.Abort()
		return
	}
	context.Set(KeyIDPlaceholder, value)
	context.Next()
}

// ValidateKey validates that `Authorization: Key` header has been provided and has correct value.
func ValidateKey(context *gin.Context, keyRepository dao.KeysRepositoryProvider) {
	value := strings.Replace(context.GetHeader("Authorization"), "Key ", "", -1)
	if value == "" {
		errors.SetHTTPError(context, errors.UserAuthorizationError)
		context.Abort()
		return
	}

	key, err := keyRepository.GetKeyByValue(value)
	if err != nil {
		errors.SetHTTPError(context, errors.InternalServerError.WithError(err))
		context.Abort()
		return
	}

	if key == nil || key.Value != value {
		errors.SetHTTPError(context, errors.UserAuthorizationError)
		context.Abort()
		return
	}
	context.Next()
}

// ValidateAdminKey validates that `Authorization: Admin Key` header has been provided and has correct value.
func ValidateAdminKey(context *gin.Context, adminKey string) {
	value := strings.Replace(context.GetHeader("Authorization"), "AdminKey ", "", -1)
	if value == "" {
		errors.SetHTTPError(context, errors.UserAuthorizationError)
		context.Abort()
		return
	}

	if value != adminKey {
		errors.SetHTTPError(context, errors.UserAuthorizationError)
		context.Abort()
		return
	}
	context.Next()
}
