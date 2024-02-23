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

var upgrader = websocket.Upgrader{}

func handlerWaitForever(t *testing.T, wg *sync.WaitGroup, letgo <-chan struct{}) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()

		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)

		<-letgo

		require.NoError(t, c.Close())
	})
}

func handlerSender(t *testing.T, wg *sync.WaitGroup, message string) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()

		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)

		err = c.WriteMessage(websocket.TextMessage, []byte(message))
		require.NoError(t, err)

		require.NoError(t, c.Close())
	})
}

func handlerReceiver(t *testing.T, wg *sync.WaitGroup, expect string) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()

		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)

		var message string
		err = c.ReadJSON(&message)
		require.NoError(t, err)

		require.Equal(t, expect, message)

		require.NoError(t, c.Close())
	})
}

func Test_Connection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("connection receive", func(t *testing.T) {
		const testMessage = "test-message"

		var wg sync.WaitGroup
		s := httptest.NewServer(handlerSender(t, &wg, testMessage))
		defer s.Close()

		address := strings.TrimPrefix(s.URL, "http://")
		conn, err := connect.NewConnector(config.Target{
			Address: address,
			WSS:     false,
		}, testClientName).Connect(ctx, gameId)
		require.NoError(t, err)
		require.NotNil(t, conn)

		data, err := conn.Receive(ctx)
		require.NoError(t, err)
		require.Equal(t, testMessage, string(data))

		require.NoError(t, conn.Close(ctx))

		wg.Wait()
	})

	t.Run("connection send", func(t *testing.T) {
		const testMessage = "another-test-message"

		var wg sync.WaitGroup
		s := httptest.NewServer(handlerReceiver(t, &wg, testMessage))
		defer s.Close()

		address := strings.TrimPrefix(s.URL, "http://")
		conn, err := connect.NewConnector(config.Target{
			Address: address,
			WSS:     false,
		}, testClientName).Connect(ctx, gameId)
		require.NoError(t, err)
		require.NotNil(t, conn)

		err = conn.Send(ctx, testMessage)
		require.NoError(t, err)

		require.NoError(t, conn.Close(ctx))

		wg.Wait()
	})

	t.Run("cancel context", func(t *testing.T) {
		letgo := make(chan struct{})

		var wg sync.WaitGroup
		s := httptest.NewServer(handlerWaitForever(t, &wg, letgo))
		defer s.Close()

		address := strings.TrimPrefix(s.URL, "http://")
		conn, err := connect.NewConnector(config.Target{
			Address: address,
			WSS:     false,
		}, testClientName).Connect(ctx, gameId)
		require.NoError(t, err)
		require.NotNil(t, conn)

		t.Run("receive", func(t *testing.T) {
			ctx1, cancel := context.WithCancel(ctx)

			cancel()

			_, err := conn.Receive(ctx1)
			require.ErrorIs(t, err, context.Canceled)
		})

		t.Run("send", func(t *testing.T) {
			ctx1, cancel := context.WithCancel(ctx)

			cancel()

			err = conn.Send(ctx1, "test-message")
			require.ErrorIs(t, err, context.Canceled)
		})

		close(letgo)

		require.NoError(t, conn.Close(ctx))

		wg.Wait()
	})
}
