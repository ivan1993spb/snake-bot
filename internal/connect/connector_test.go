package connect_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-bot/internal/config"
	"github.com/ivan1993spb/snake-bot/internal/connect"
)

const testClientName = "test-client-name"

const gameId = 1

func Test_Connector_FailedConnection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Do not upgrade connection
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	address := strings.TrimPrefix(s.URL, "http://")
	connector := connect.NewConnector(config.Target{
		Address: address,
		WSS:     false,
	}, testClientName)

	t.Run("connect failed", func(t *testing.T) {
		conn, err := connector.Connect(ctx, gameId)
		require.Error(t, err)
		require.Nil(t, conn)
	})
}

func Test_Connector_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var upgrader = websocket.Upgrader{}

	var wg sync.WaitGroup

	s := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wg.Add(1)
			defer wg.Done()

			c, err := upgrader.Upgrade(w, r, nil)
			require.NoError(t, err)

			clientName := r.Header.Get(connect.HeaderClientName)
			require.Equal(t, testClientName, clientName)

			require.NoError(t, c.Close())
		},
	))
	defer s.Close()

	address := strings.TrimPrefix(s.URL, "http://")
	connector := connect.NewConnector(config.Target{
		Address: address,
		WSS:     false,
	}, testClientName)

	t.Run("connect success", func(t *testing.T) {
		conn, err := connector.Connect(ctx, gameId)
		require.NoError(t, err)
		require.NotNil(t, conn)
		require.NoError(t, conn.Close(ctx))
	})

	wg.Wait()
}
