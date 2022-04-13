package custom

import (
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// LogLevel represents LogLevel config parameter.
type LogLevel log.Level

// Decode decodes LogLevel ENV parameter.
func (l *LogLevel) Decode(value string) error {

	logLevel, err := log.ParseLevel(strings.ToLower(value))
	if err != nil {
		return errors.New("unsupported `LOG_LEVEL` parameter value")
	}

	*l = LogLevel(logLevel)

	return nil
}
