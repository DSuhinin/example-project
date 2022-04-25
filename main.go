package main

import (
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	"github.com/DSuhinin/passbase-test-task/core"

	"github.com/DSuhinin/passbase-test-task/app"
	"github.com/DSuhinin/passbase-test-task/app/config"
	"github.com/DSuhinin/passbase-test-task/app/controller"
	"github.com/DSuhinin/passbase-test-task/app/service/currencies"
	"github.com/DSuhinin/passbase-test-task/app/service/currencies/fixer"
	"github.com/DSuhinin/passbase-test-task/app/service/keys"
	"github.com/DSuhinin/passbase-test-task/app/service/keys/dao"
)

// main entry point
func main() {

	// 1. setup config.
	appCfg, err := config.New()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("error initializing config")
	}

	// 2. init logger.
	core.InitJSONLogger(log.Level(appCfg.LogLevel))

	// 3. initialize db connections.
	dbConnection, err := core.NewDB().GetConnection(
		appCfg.DatabaseUser,
		appCfg.DatabasePass,
		core.PostgresType,
		appCfg.DatabaseName,
		appCfg.DatabaseHost,
	)
	if err != nil {
		log.Fatalf("error establishing connection to %s database: %+v", core.PostgresType, err)
	}

	// 4. init cache
	cache := cache.New(5*time.Minute, 10*time.Minute)

	// 5. create main router and run service.
	router, err := app.NewRouter(appCfg, controller.New(
		keys.NewService(
			dao.NewKeysRepository(dbConnection),
		),
		currencies.NewService(
			fixer.NewClient(appCfg.FixerAPIBaseURL, appCfg.FixerAPIKey, cache),
		),
	), dao.NewKeysRepository(dbConnection))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("error configuring `router")
	}

	if err := router.Start(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("error running service")
	}
}
