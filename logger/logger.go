package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// InitLogger initializes the logger with desired configuration
func InitLogger() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	return log
}
