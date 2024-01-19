package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

func NewRequestLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(RequestLogger{})
}

type RequestLogger struct{}

func (RequestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	ctx := r.Context()
	log := utils.GetLogger(ctx)

	logFields := logrus.Fields{}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method
	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()
	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry := &StructuredLoggerEntry{
		Log: log.WithFields(logFields),
	}

	entry.Log.Infoln("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Log *logrus.Entry
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header,
	elapsed time.Duration, extra interface{}) {
	l.Log = l.Log.WithFields(logrus.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed_ms":   float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Log.Infoln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Log = l.Log.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
