package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DSuhinin/passbase-test-task/core/errors"

	"github.com/DSuhinin/passbase-test-task/app/api/response"
	"github.com/DSuhinin/passbase-test-task/app/middleware"
)

// CreateKey handles `POST /keys` route.
func (c Controller) CreateKey(ctx *gin.Context) {
	key, err := c.keysService.CreateKey()
	if err != nil {
		errors.SetHTTPError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.NewKey(key))
}

// GetKeys handles `GET /keys` route.
func (c Controller) GetKeys(ctx *gin.Context) {
	keys, err := c.keysService.GetKeys()
	if err != nil {
		errors.SetHTTPError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewKeys(keys))
}

// RegenerateKey handles `PUT /keys/:key_id/regenerate` route.
func (c Controller) RegenerateKey(ctx *gin.Context) {
	key, err := c.keysService.RegenerateKey(ctx.GetInt(middleware.KeyIDPlaceholder))
	if err != nil {
		errors.SetHTTPError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewKey(key))
}

// DeleteKey handles `DELETE /keys/:key_id` route.
func (c Controller) DeleteKey(ctx *gin.Context) {
	if err := c.keysService.DeleteKey(ctx.GetInt(middleware.KeyIDPlaceholder)); err != nil {
		errors.SetHTTPError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
