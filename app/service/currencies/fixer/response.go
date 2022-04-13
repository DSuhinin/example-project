package fixer

// CurrencyExchangeResponse represents Response object for `http://data.fixer.io/api/latest` endpoint.
type CurrencyExchangeResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     struct {
		Usd float64 `json:"USD"`
	} `json:"rates"`
}
