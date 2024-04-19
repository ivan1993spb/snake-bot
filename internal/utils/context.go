package utils

import (
	"context"
	"path"
	"strconv"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

type (
	loggerKey struct{}
	moduleKey struct{}
	taskKey   struct{}
)

func WithLogger(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

func GetLogger(ctx context.Context) *logrus.Entry {
	if log, ok := ctx.Value(loggerKey{}).(*logrus.Entry); ok {
		return log
	}
	return DiscardEntry
}

func WithFields(ctx context.Context, fields logrus.Fields) context.Context {
	log := GetLogger(ctx).WithFields(fields)
	return WithLogger(ctx, log)
}

func WithField(ctx context.Context, key string, value interface{}) context.Context {
	log := GetLogger(ctx).WithField(key, value)
	return WithLogger(ctx, log)
}

func WithModule(ctx context.Context, module string) context.Context {
	parent := GetModulePath(ctx)

	if parent != "" {
		// don't re-append module when module is the same.
		if path.Base(parent) == module {
			return ctx
		}

		module = path.Join(parent, module)
	}

	ctx = context.WithValue(ctx, moduleKey{}, module)
	return WithField(ctx, "module", module)
}

func GetModulePath(ctx context.Context) string {
	if module, ok := ctx.Value(moduleKey{}).(string); ok {
		return module
	}

	return ""
}

var taskIdCounter uint64

func GetTaskId(ctx context.Context) uint64 {
	if taskId, ok := ctx.Value(taskKey{}).(uint64); ok {
		return taskId
	}

	return 0
}

func WithTaskId(ctx context.Context) context.Context {
	if taskId := GetTaskId(ctx); taskId != 0 {
		return ctx
	}

	taskId := atomic.AddUint64(&taskIdCounter, 1)

	ctx = context.WithValue(ctx, taskKey{}, taskId)
	return WithField(ctx, "task", strconv.FormatUint(taskId, 10))
}
