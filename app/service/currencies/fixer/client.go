package fixer

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DSuhinin/passbase-test-task/core/errors"
)

// ClientProvider provides and interface to work with Fixer service.
type ClientProvider interface {
	// GetExchangeRate returns current exchange rate.
	GetExchangeRate() (float64, error)
}

// Client represents HTTP client to work with Fixer service.
type Client struct {
	cache   *cache.Cache
	token   string
	baseURL string
}

// NewClient creates new HTTP client to work with Fixer service.
func NewClient(baseURL, token string, cache *cache.Cache) *Client {
	return &Client{
		cache:   cache,
		token:   token,
		baseURL: baseURL,
	}
}

// GetExchangeRate returns current exchange rate.
func (c Client) GetExchangeRate() (float64, error) {
	if c.cache != nil {
		rate, ok := c.cache.Get("exchange_rate")
		if ok {
			return rate.(float64), nil
		}
	}

	//nolint
	resp, err := http.Get(
		fmt.Sprintf("%s/latest?access_key=%s&base=EUR&symbols=USD", c.baseURL, c.token),
	)
	if err != nil {
		return 0, errors.Wrap(err, "error getting exchange information")
	}

	if resp == nil {
		return 0, errors.Wrap(err, "fixer service returned an empty response")
	}
	//nolint
	defer resp.Body.Close()

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

	if c.cache != nil {
		c.cache.Set("exchange_rate", data.Rates.Usd, 10*time.Second)
	}

	return data.Rates.Usd, nil

}
