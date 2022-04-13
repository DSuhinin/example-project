package currencies

import (
	"github.com/DSuhinin/passbase-test-task/core/errors"

	"github.com/DSuhinin/passbase-test-task/app/api/request"
	"github.com/DSuhinin/passbase-test-task/app/service/currencies/fixer"
)

// ServiceProvider provides Currencies related operations.
type ServiceProvider interface {
	// CurrenciesExchange make currency exchange.
	CurrenciesExchange(req *request.CurrencyExchange) (float64, error)
}

// Service implements ServiceProvider interface.
type Service struct {
	fixerService fixer.ClientProvider
}

// NewService creates new Currencies service instance.
func NewService(fixerService fixer.ClientProvider) *Service {
	return &Service{
		fixerService: fixerService,
	}
}

// CurrenciesExchange make currency exchange.
func (s Service) CurrenciesExchange(req *request.CurrencyExchange) (float64, error) {
	if err := ValidateCurrenciesExchangeRequest(req); err != nil {
		return 0, err
	}

	rate, err := s.fixerService.GetExchangeRate()
	if err != nil {
		return 0, errors.InternalServerError.WithError(err)
	}

	if req.From == "EUR" {
		return rate * req.Amount, nil
	}

	return (1 / rate) * req.Amount, nil
}
