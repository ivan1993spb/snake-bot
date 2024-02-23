package connect

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/ivan1993spb/snake-bot/internal/config"
)

type Connector struct {
	address string
	wss     bool

	dialer *websocket.Dialer
	header http.Header
}

const handshakeTimeout = 45 * time.Second

const HeaderClientName = "X-Snake-Client"

func NewConnector(cfg config.Target, clientName string) *Connector {
	c := &Connector{
		address: cfg.Address,
		wss:     cfg.WSS,
	}

	c.dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: handshakeTimeout,
	}

	c.header = http.Header{}
	c.header.Add(HeaderClientName, clientName)

	return c
}

const gameWebsocketPathFormat = "/ws/games/%d"

func getGameWebSocketPath(gameId int) string {
	return fmt.Sprintf(gameWebsocketPathFormat, gameId)
}

func (c *Connector) getGameWebWocketURL(gameId int) *url.URL {
	u := &url.URL{
		Host:   c.address,
		Scheme: "ws",
	}

	if c.wss {
		u.Scheme = "wss"
	}

	u.Path = getGameWebSocketPath(gameId)

	return u
}

func (c *Connector) Connect(ctx context.Context, gameId int) (Connection, error) {
	u := c.getGameWebWocketURL(gameId).String()
	ws, resp, err := c.dialer.DialContext(ctx, u, c.header)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

	conn := &connection{
		conn: ws,
		resp: resp,
	}

	return conn, nil
}
