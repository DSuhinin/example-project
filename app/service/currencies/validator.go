package currencies

import (
	"fmt"
	"github.com/DSuhinin/passbase-test-task/app/api/request"
	"github.com/DSuhinin/passbase-test-task/core/errors"
)

// Currency constants
const (
	USDCurrency = "USD"
	EURCurrency = "EUR"
)

// AllowedCurrencyList supported list of Currencies.
var (
	AllowedCurrencyList = map[string]bool{
		USDCurrency: true,
		EURCurrency: true,
	}
)

// ValidateCurrenciesExchangeRequest validates `GET /currencies/exchange` endpoint request.
func ValidateCurrenciesExchangeRequest(req *request.CurrencyExchange) error {
	if _, ok := AllowedCurrencyList[req.From]; !ok {
		return errors.NewHTTPBadRequest(
			20000, fmt.Sprintf("`from` parameter should be %s or %s", USDCurrency, EURCurrency),
		)
	}

	if _, ok := AllowedCurrencyList[req.To]; !ok {
		return errors.NewHTTPBadRequest(
			20010, fmt.Sprintf("`to` parameter should be %s or %s", USDCurrency, EURCurrency),
		)
	}

	if req.Amount <= 0 {
		return errors.NewHTTPBadRequest(20020, "`amount` parameter should be grater then zero")
	}

	return nil
}
