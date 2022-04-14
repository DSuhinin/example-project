package response

import "math"

// CurrencyExchange represents Response object for `GET /currencies/exchange` endpoint.
type CurrencyExchange struct {
	Result float64 `json:"result"`
}

// NewCurrencyExchange creates new instance of Response object for `GET /currencies/exchange` endpoint.
func NewCurrencyExchange(result float64) *CurrencyExchange {
	return &CurrencyExchange{
		Result: math.Round(result*100) / 100,
	}
}
