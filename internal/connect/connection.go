package connect

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var _ Connection = (*connection)(nil)

type Connection interface {
	Send(data interface{}) error
	Receive() ([]byte, error)
	Close() error
}

type connection struct {
	conn *websocket.Conn
	resp *http.Response
	mux  *sync.Mutex
}

const deadlinePing = time.Second

func (c *connection) Send(data interface{}) error {
	// TODO: Ping?
	//deadline := time.Now().Add(deadlinePing)
	// TODO: Handle error.
	//c.conn.WriteControl(websocket.PingMessage, nil, deadline)
	// TODO: Need concurrency!
	//c.mux.Lock()
	//defer c.mux.Unlock()
	err := c.conn.WriteJSON(data)
	if err != nil {
		return errors.Wrap(err, "sending data")
	}
	return nil
}

var ErrConnectionClosed = errors.New("connection has been closed")

var ErrWrongMsgType = errors.New("wrong message type")

func (c *connection) Receive() ([]byte, error) {
	// TODO: Need concurrency!
	//c.mux.Lock()
	//defer c.mux.Unlock()
	messageType, message, err := c.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			return nil, ErrConnectionClosed
		}
		return nil, errors.Wrap(err, "reading message")
	}
	if messageType != websocket.TextMessage {
		return nil, ErrWrongMsgType
	}
	return message, nil
}

func (c *connection) Close() error {
	// TODO: Need concurrency!
	//c.mux.Lock()
	//defer c.mux.Unlock()
	{
		message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
		err := c.conn.WriteMessage(websocket.CloseMessage, message)
		if err != nil {
			return errors.Wrap(err, "writing close message")
		}
	}
	err := c.conn.Close()
	if err != nil {
		return errors.Wrap(err, "closing websocket")
	}
	return nil
}
