package controller

import (
	"github.com/DSuhinin/passbase-test-task/app/service/currencies"
	"github.com/DSuhinin/passbase-test-task/app/service/keys"
)

// Controller represents Controller layer to handle incoming requests.
type Controller struct {
	keysService       keys.ServiceProvider
	currenciesService currencies.ServiceProvider
}

// New creates new instance of Controller.
func New(keysService keys.ServiceProvider, currenciesService currencies.ServiceProvider) *Controller {
	return &Controller{
		keysService:       keysService,
		currenciesService: currenciesService,
	}
}
