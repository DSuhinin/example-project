package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/DSuhinin/passbase-test-task/app/config/custom"
)

// Config represents Lambda Config object.
type Config struct {
	LogLevel        custom.LogLevel `envconfig:"LOG_LEVEL" default:"error"`
	AdminKey        string          `envconfig:"ADMIN_KEY" default:"supersecurekey" required:"true"`
	FixerAPIKey     string          `envconfig:"FIXER_API_KEY" required:"true"`
	FixerAPIBaseURL string          `envconfig:"FIXER_API_BASE_URL" required:"true"`
	DatabaseUser    string          `envconfig:"DATABASE_USER" required:"true"`
	DatabasePass    string          `envconfig:"DATABASE_PASS" required:"true"`
	DatabaseName    string          `envconfig:"DATABASE_NAME" required:"true"`
	DatabaseHost    string          `envconfig:"DATABASE_HOST" required:"true"`
	ServerAddress   string          `envconfig:"SERVER_ADDRESS" default:"127.0.0.1:8080" required:"true"`
}

// New creates new instance of Config object.
func New() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, errors.Wrap(err, "error creating config from ENV")
	}

	return &c, nil
}
