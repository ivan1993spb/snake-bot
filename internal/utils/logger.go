package utils

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-bot/internal/config"
)

const (
	loggerDisableColors   = true
	loggerTimestampFormat = time.RFC1123
)

var DiscardLogger = &logrus.Logger{
	Out:          io.Discard,
	Hooks:        make(logrus.LevelHooks),
	Formatter:    new(NullFormatter),
	Level:        logrus.PanicLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
}

var DiscardEntry = logrus.NewEntry(DiscardLogger)

type NullFormatter struct{}

func (*NullFormatter) Format(*logrus.Entry) ([]byte, error) {
	return nil, nil
}

const defaultLogLevel = logrus.InfoLevel

func NewLogger(cfg config.Log) *logrus.Entry {
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
		logger.SetLevel(defaultLogLevel)
	} else {
		logger.SetLevel(level)
	}

	return logrus.NewEntry(logger)
}
