package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DSuhinin/passbase-test-task/core/errors"

	"github.com/DSuhinin/passbase-test-task/app/api/request"
	"github.com/DSuhinin/passbase-test-task/app/api/response"
)

// CurrenciesExchange handles `GET /currencies/exchange` route.
func (c Controller) CurrenciesExchange(ctx *gin.Context) {
	var req request.CurrencyExchange
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errors.SetHTTPError(ctx, errors.QueryParametersParsingError.WithError(
			errors.Wrap(err, "impossible to parse request parameters"),
		))
		return
	}
	amount, err := c.currenciesService.CurrenciesExchange(&req)
	if err != nil {
		errors.SetHTTPError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewCurrencyExchange(amount))
}
