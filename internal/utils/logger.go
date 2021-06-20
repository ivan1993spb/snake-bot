package utils

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-bot/internal/config"
)

const (
	loggerDisableColors   = true
	loggerTimestampFormat = time.RFC1123
)

func NewLogger(cfg config.Log) *logrus.Logger {
	logger := logrus.New()
	if cfg.EnableJSON {
		logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: loggerTimestampFormat,
		}
	} else {
		logger.Formatter = &logrus.TextFormatter{
			DisableColors:   loggerDisableColors,
			TimestampFormat: loggerTimestampFormat,
		}
	}
	if level, err := logrus.ParseLevel(cfg.Level); err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}
	return logger
}
