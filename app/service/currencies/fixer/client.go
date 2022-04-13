package fixer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DSuhinin/passbase-test-task/core/errors"
)

// ClientProvider provides and interface to work with Fixer service.
type ClientProvider interface {
	// GetExchangeRate returns current exchange rate.
	GetExchangeRate() (float64, error)
}

// Client represents HTTP client to work with Fixer service.
type Client struct {
	token   string
	baseURL string
}

// NewClient creates new HTTP client to work with Fixer service.
func NewClient(baseURL, token string) *Client {
	return &Client{
		token:   token,
		baseURL: baseURL,
	}
}

// GetExchangeRate returns current exchange rate.
func (c Client) GetExchangeRate() (float64, error) {
	resp, err := http.Get(
		fmt.Sprintf("%s/latest?access_key=%s&base=EUR&symbols=USD", c.baseURL, c.token),
	)
	if err != nil {
		return 0, errors.Wrap(err, "error getting exchange information")
	}

	if resp == nil {
		return 0, errors.Wrap(err, "fixer service returned empty response")
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errors.Errorf("fixer service returned non 200 status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrap(err, "error reading response body")
	}

	data := CurrencyExchangeResponse{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, errors.Wrap(err, "error unmarshaling json data")
	}

	return data.Rates.Usd, nil

}
