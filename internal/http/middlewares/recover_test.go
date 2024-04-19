package middlewares_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-bot/internal/http/middlewares"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type panicHandler struct {
	mux sync.Mutex

	file string
	line int
	ok   bool

	count int
	msg   string
}

func (h *panicHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.count++

	_, h.file, h.line, h.ok = runtime.Caller(0)
	panic(h.msg)
}

func TestRecover(t *testing.T) {
	const panicMessage = "this is expected panic"

	logger, hook := test.NewNullLogger()
	ctx := context.Background()
	ctx = utils.WithLogger(ctx, logrus.NewEntry(logger))

	r := chi.NewRouter()
	r.Use(middlewares.RequestID)
	r.Use(middlewares.NewRequestLogger())
	r.Use(middleware.Recoverer)
	handler := &panicHandler{
		msg: panicMessage,
	}
	r.Method("GET", "/", handler)

	server := httptest.NewUnstartedServer(r)
	server.Config.BaseContext = func(net.Listener) context.Context {
		return utils.WithModule(ctx, "test")
	}

	server.StartTLS()
	defer server.Close()

	resp, err := server.Client().Get(server.URL)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.True(t, handler.ok)
	require.Equal(t, 1, handler.count)

	require.Equal(t, 500, resp.StatusCode)
	content, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Empty(t, content)

	entries := hook.AllEntries()
	require.Len(t, entries, 3)
	require.Equal(t, "request started", entries[0].Message)
	require.Equal(t, "recovered panic", entries[1].Message)
	require.Equal(t, panicMessage, entries[1].Data["panic"])
	require.NotNil(t, entries[1].Data["pc"])
	// Check the next line: +1
	require.Equal(t, handler.line+1, entries[1].Data["line"])
	require.Equal(t, handler.file, entries[1].Data["file"])
	require.Equal(t, "request complete", entries[2].Message)
}
