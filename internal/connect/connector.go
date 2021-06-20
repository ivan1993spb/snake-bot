package connect

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ivan1993spb/snake-bot/internal/config"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	websocketInsecureScheme = "ws"
	websocketSecureScheme   = "wss"
)

const gameWebsocketPathFormat = "/ws/games/%d"

const handshakeTimeout = 45 * time.Second

const headerClientName = "X-Snake-Client"

type Connector struct {
	address string
	wss     bool

	dialer *websocket.Dialer
	http.Header
}

func NewConnector(cfg config.Target, clientName string) *Connector {
	c := &Connector{
		address: cfg.Address,
		wss:     cfg.WSS,
	}

	c.dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: handshakeTimeout,
	}

	c.Header = http.Header{}
	c.Header.Add(headerClientName, clientName)

	return c
}

func getGameWebSocketPath(gameId int) string {
	return fmt.Sprintf(gameWebsocketPathFormat, gameId)
}

func (c *Connector) getGameWebWocketURL(gameId int) *url.URL {
	u := &url.URL{
		Host:   c.address,
		Scheme: websocketInsecureScheme,
	}

	if c.wss {
		u.Scheme = websocketSecureScheme
	}

	u.Path = getGameWebSocketPath(gameId)

	return u
}

func (c *Connector) Connect(ctx context.Context, gameId int) (Connection, error) {
	u := c.getGameWebWocketURL(gameId)
	ws, resp, err := c.dialer.DialContext(ctx, u.String(), c.Header)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect")
	}
	conn := &connection{
		conn: ws,
		resp: resp,
		mux:  &sync.Mutex{},
	}
	return conn, nil
}
