package connect

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var _ Connection = (*connection)(nil)

//counterfeiter:generate . Connection
type Connection interface {
	Send(ctx context.Context, data interface{}) error
	Receive(ctx context.Context) ([]byte, error)
	Close(ctx context.Context) error
}

type connection struct {
	conn *websocket.Conn
	resp *http.Response

	// gorilla websocket supports 1 reader and 1 writer
	// concurrently
	rmux sync.Mutex
	wmux sync.Mutex
}

func (c *connection) Send(ctx context.Context, data interface{}) error {
	errCh := make(chan error, 1)

	go func() {
		c.wmux.Lock()
		defer c.wmux.Unlock()

		defer close(errCh)

		errCh <- c.conn.WriteJSON(data)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, "sending data")
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

var ErrConnectionClosed = errors.New("connection has been closed")

var ErrWrongMsgType = errors.New("wrong message type")

type readResult struct {
	messageType int
	message     []byte
	err         error
}

func (c *connection) Receive(ctx context.Context) ([]byte, error) {
	ch := make(chan readResult, 1)

	go func() {
		c.rmux.Lock()
		defer c.rmux.Unlock()

		defer close(ch)

		messageType, message, err := c.conn.ReadMessage()
		ch <- readResult{
			messageType: messageType,
			message:     message,
			err:         err,
		}
	}()

	select {
	case result := <-ch:
		messageType := result.messageType
		message := result.message
		err := result.err

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil, ErrConnectionClosed
			}
			return nil, errors.Wrap(err, "reading message")
		}

		if messageType != websocket.TextMessage {
			return nil, ErrWrongMsgType
		}

		return message, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *connection) Close(ctx context.Context) error {
	// If the context is canceled, close the connection immediately.
	if errors.Is(ctx.Err(), context.Canceled) {
		return c.conn.Close()
	}

	errCh := make(chan error, 1)

	go func() {
		c.wmux.Lock()
		defer c.wmux.Unlock()

		defer close(errCh)

		errCh <- c.conn.Close()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, "closing websocket")
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
