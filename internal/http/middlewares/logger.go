package middlewares

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

func NewRequestLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(RequestLogger{})
}

type RequestLogger struct{}

// NewLogEntry creates an alternative log entry which is stored in
// the context with the key middleware.LogEntryCtxKey. This log entry
// is used inside chi's own middlewares: Recoverer, e.g.
//
// The new log entry has some new fields which aren't really needed
// in handlers although could be used in the future.
func (RequestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	ctx := r.Context()

	// The logger from the context is supposed to already have the
	// request id providing that the log middleware was added
	// after RequestID.
	log := utils.GetLogger(ctx)

	fields := logrus.Fields{}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	fields["http_scheme"] = scheme
	fields["http_proto"] = r.Proto
	fields["http_method"] = r.Method
	fields["remote_addr"] = r.RemoteAddr
	fields["user_agent"] = r.UserAgent()
	fields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry := &StructuredLoggerEntry{
		// There is no need to save the fields: the log has the request id
		Logger: log,
	}

	entry.Logger.WithFields(fields).Infoln("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger *logrus.Entry
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header,
	elapsed time.Duration, extra interface{}) {

	fields := logrus.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed_ms":   float64(elapsed.Nanoseconds()) / 1000000.0,
	}

	l.Logger.WithFields(fields).Infoln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	fields := logrus.Fields{
		"panic": fmt.Sprintf("%+v", v),
	}

	// Skips: 0=current, 1=recoverer, 2=panic
	pc, file, line, ok := runtime.Caller(3)
	if ok {
		fields["pc"] = pc
		fields["file"] = file
		fields["line"] = line
	}

	l.Logger.WithFields(fields).Error("recovered panic")
}
