package utils

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey struct{}

var loggerKey = contextKey{}

func LogContext(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

var defaultLog = logrus.New()

func Log(ctx context.Context) *logrus.Logger {
	if logger, ok := ctx.Value(loggerKey).(*logrus.Logger); ok {
		return logger
	}
	return defaultLog
}
