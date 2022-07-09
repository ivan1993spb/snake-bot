package utils

import (
	"context"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type contextKey struct{}

var loggerKey = contextKey{}

func LogContext(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

var defaultLog = &Logger{
	Out:          io.Discard,
	Formatter:    new(NullFormatter),
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.PanicLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
}

func Log(ctx context.Context) *logrus.Logger {
	// TODO: Return *logrus.Entry.
	if logger, ok := ctx.Value(loggerKey).(*logrus.Logger); ok {
		return logger
	}
	return defaultLog
}
