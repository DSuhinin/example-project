package core

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// InitJSONLogger initializes Logrus to log to stdout in JSON format, beginning with time, level and message fields.
func InitJSONLogger(logLevel log.Level) {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: log.FieldMap{ // assure field ordering
			log.FieldKeyTime:  "@@time",
			log.FieldKeyLevel: "@level",
			log.FieldKeyMsg:   "@message",
		},
	})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Set log level
	log.SetLevel(logLevel)
}
