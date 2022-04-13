package request

// CurrencyExchange represents Request object for `GET /currencies/exchange` endpoint.
type CurrencyExchange struct {
	To     string  `form:"to"`
	From   string  `form:"from"`
	Amount float64 `form:"amount"`
}
