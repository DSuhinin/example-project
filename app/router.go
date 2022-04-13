package app

import (
	"github.com/DSuhinin/passbase-test-task/core/errors"
	"github.com/gin-gonic/gin"

	"github.com/DSuhinin/passbase-test-task/app/config"
	"github.com/DSuhinin/passbase-test-task/app/controller"
	"github.com/DSuhinin/passbase-test-task/app/middleware"
	"github.com/DSuhinin/passbase-test-task/app/service/keys/dao"
)

// supported route list.
//nolint:lll
const (
	CreateKeyRoute          = "/keys"
	GetKeysRoute            = "/keys"
	RegenerateKeyRoute      = "/keys/:key_id/regenerate"
	DeleteKeyRoute          = "/keys/:key_id"
	CurrenciesExchangeRoute = "/currencies/exchange"
)

// Router represents Application router object.
type Router struct {
	config    *config.Config
	ginEngine *gin.Engine
}

// NewRouter initializes the gin router and routes.
func NewRouter(
	config *config.Config,
	controller *controller.Controller,
	keysRepository dao.KeysRepositoryProvider,
) (*Router, error) {
	g := gin.New()
	g.Use(gin.Recovery())

	if err := g.SetTrustedProxies(nil); err != nil {
		return nil, errors.Wrap(err, "error configuring router")
	}

	g.POST(
		CreateKeyRoute,
		func(context *gin.Context) {
			middleware.ValidateAdminKey(context, config.AdminKey)
		},
		controller.CreateKey,
	)
	g.GET(
		GetKeysRoute,
		func(context *gin.Context) {
			middleware.ValidateAdminKey(context, config.AdminKey)
		},
		controller.GetKeys,
	)
	g.PUT(
		RegenerateKeyRoute,
		func(context *gin.Context) {
			middleware.ValidateAdminKey(context, config.AdminKey)
		},
		middleware.ValidateKeyID,
		controller.RegenerateKey,
	)
	g.DELETE(
		DeleteKeyRoute,
		func(context *gin.Context) {
			middleware.ValidateAdminKey(context, config.AdminKey)
		},
		middleware.ValidateKeyID,
		controller.DeleteKey,
	)
	g.GET(
		CurrenciesExchangeRoute,
		func(context *gin.Context) {
			middleware.ValidateKey(context, keysRepository)
		},
		controller.CurrenciesExchange,
	)

	return &Router{
		config:    config,
		ginEngine: g,
	}, nil
}

// Start starts application routing.
func (r Router) Start() error {
	return r.ginEngine.Run(r.config.ServerAddress)
}
