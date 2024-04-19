package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/ivan1993spb/snake-bot/internal/connect"
	"github.com/ivan1993spb/snake-bot/internal/connect/connectfakes"
	"github.com/ivan1993spb/snake-bot/internal/core"
	"github.com/ivan1993spb/snake-bot/internal/core/corefakes"
	"github.com/ivan1993spb/snake-bot/internal/types"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type receiveReturns struct {
	data []byte
	err  error
}

func connectionStub(
	t *testing.T,
	receiveCh <-chan *receiveReturns,
	sendCh <-chan error,
	closeCh <-chan error,
) connect.Connection {

	t.Helper()

	connection := &connectfakes.FakeConnection{
		ReceiveStub: func(ctx context.Context) ([]byte, error) {
			t.Helper()

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case ret, ok := <-receiveCh:
				if !ok {
					return nil, errors.New("receive channel is closed")
				}
				return ret.data, ret.err
			}
		},
		SendStub: func(ctx context.Context, data interface{}) error {
			t.Helper()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-sendCh:
				return err
			}
		},
		CloseStub: func(ctx context.Context) error {
			t.Helper()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case resErr := <-closeCh:
				return resErr
			}
		},
	}

	return connection
}

type connectReturns struct {
	connection connect.Connection
	err        error
}

func connectorStub(t *testing.T, connectCh <-chan *connectReturns) core.Connector {
	t.Helper()

	connector := &corefakes.FakeConnector{
		ConnectStub: func(ctx context.Context, gameId int) (connect.Connection, error) {
			t.Helper()

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case ret, ok := <-connectCh:
				if !ok {
					return nil, errors.New("connect channel is closed")
				}
				return ret.connection, ret.err
			}
		},
	}

	return connector
}

type runReturns struct {
	ch chan types.Direction
}

func botEngineStub(t *testing.T, runCh <-chan *runReturns) core.BotEngine {
	t.Helper()

	botEngine := &corefakes.FakeBotEngine{
		RunStub: func(ctx context.Context) <-chan types.Direction {
			t.Helper()

			select {
			case <-ctx.Done():
				return nil
			case ret, ok := <-runCh:
				if !ok {
					return nil
				}

				chout := make(chan types.Direction)

				go func() {
					defer close(chout)

					for {
						select {
						case <-ctx.Done():
							return
						case dir, ok := <-ret.ch:
							if !ok {
								return
							}

							select {
							case <-ctx.Done():
								return
							case chout <- dir:
							}
						}
					}
				}()

				return chout
			}
		},
	}

	return botEngine
}

func parserStub(t *testing.T, parseCh <-chan error) core.Parser {
	t.Helper()

	parser := &corefakes.FakeParser{
		ParseStub: func(message []byte) error {
			return <-parseCh
		},
	}

	return parser
}

var directions = []types.Direction{
	types.DirectionSouth,
	types.DirectionWest,
	types.DirectionNorth,
	types.DirectionEast,
}

const testBotOperatorTimeout = time.Millisecond * 100

func Test_BotOperator(t *testing.T) {
	const (
		gameId = 99
		n      = 100
	)

	ctx, cancel := context.WithTimeout(context.Background(), testBotOperatorTimeout)
	defer cancel()

	ctx = utils.WithLogger(ctx, utils.DiscardEntry)

	fakeRand := &corefakes.FakeRand{}
	clock := utils.ImmediatelyClock
	connectCh := make(chan *connectReturns)
	defer close(connectCh)
	connector := connectorStub(t, connectCh)
	runCh := make(chan *runReturns)
	defer close(runCh)
	botEngine := botEngineStub(t, runCh)
	parseCh := make(chan error)
	defer close(parseCh)
	parser := parserStub(t, parseCh)

	botOperator := core.NewBotOperator(&core.BotOperatorParams{
		GameId:    gameId,
		Connector: connector,
		BotEngine: botEngine,
		Parser:    parser,
		Rand:      fakeRand,
		Clock:     clock,
	})

	go botOperator.Run(ctx)

	t.Run("connect", func(t *testing.T) {
		t.Run("failed", func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case connectCh <- &connectReturns{
				connection: nil,
				err:        errors.New("connection error"),
			}:
			}
		})

		t.Run("success", func(t *testing.T) {
			receiveCh := make(chan *receiveReturns)
			defer close(receiveCh)
			sendCh := make(chan error)
			defer close(sendCh)
			closeCh := make(chan error)
			defer close(closeCh)
			directionsCh := make(chan types.Direction)
			defer close(directionsCh)

			connection := connectionStub(t, receiveCh, sendCh, closeCh)

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case connectCh <- &connectReturns{
				connection: connection,
			}:
			}

			t.Run("running bot", func(t *testing.T) {
				select {
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				case runCh <- &runReturns{
					ch: directionsCh,
				}:
				}

				for i := 0; i < n; i++ {
					select {
					case <-ctx.Done():
						t.Fatal(ctx.Err())
					case directionsCh <- directions[i%len(directions)]:
					}

					select {
					case <-ctx.Done():
						t.Fatal(ctx.Err())
					case sendCh <- nil:
					}
				}

				for i := 0; i < n; i++ {
					select {
					case <-ctx.Done():
					case receiveCh <- &receiveReturns{
						data: []byte("message"),
						err:  nil,
					}:
					}

					select {
					case <-ctx.Done():
						t.Fatal(ctx.Err())
					case parseCh <- nil:
					}
				}
			})

			t.Run("break connection when sending", func(t *testing.T) {
				select {
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				case directionsCh <- types.DirectionNorth:
				}

				select {
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				case sendCh <- connect.ErrConnectionClosed:
				}
			})

			t.Run("close connection", func(t *testing.T) {
				select {
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				case closeCh <- nil:
				}
			})
		})
	})

	t.Run("try to reconnect multiple times", func(t *testing.T) {
		const reconnectAttempts = 20

		for i := 0; i < reconnectAttempts; i++ {
			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case connectCh <- &connectReturns{
				connection: nil,
				err:        errors.New("connection error"),
			}:
			}
		}
	})

	t.Run("reconnect successfully", func(t *testing.T) {
		receiveCh := make(chan *receiveReturns, n)
		defer close(receiveCh)
		sendCh := make(chan error, n)
		defer close(sendCh)
		closeCh := make(chan error)
		defer close(closeCh)
		directionsCh := make(chan types.Direction, n)
		defer close(directionsCh)

		connection := connectionStub(t, receiveCh, sendCh, closeCh)

		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case connectCh <- &connectReturns{
			connection: connection,
		}:
		}

		t.Run("running bot again", func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case runCh <- &runReturns{
				ch: directionsCh,
			}:
			}

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case directionsCh <- types.DirectionNorth:
			}

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case sendCh <- nil:
			}

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case receiveCh <- &receiveReturns{
				data: []byte("message"),
				err:  nil,
			}:
			}

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case parseCh <- nil:
			}

			t.Run("failed to parse", func(t *testing.T) {
				select {
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				case receiveCh <- &receiveReturns{
					data: []byte("message"),
					err:  nil,
				}:
				}

				select {
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				case parseCh <- errors.New("parse error"):
				}
			})
		})

		t.Run("break connection when receiving", func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case receiveCh <- &receiveReturns{
				data: nil,
				err:  errors.New("connection error"),
			}:
			}
		})

		t.Run("close connection failure", func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case closeCh <- errors.New("connection error"):
			}
		})
	})

	t.Run("reconnect successfully again", func(t *testing.T) {
		receiveCh := make(chan *receiveReturns, n)
		defer close(receiveCh)
		sendCh := make(chan error, n)
		defer close(sendCh)
		closeCh := make(chan error)
		defer close(closeCh)
		directionsCh := make(chan types.Direction, n)
		defer close(directionsCh)

		connection := connectionStub(t, receiveCh, sendCh, closeCh)

		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case connectCh <- &connectReturns{
			connection: connection,
		}:
		}

		t.Run("running bot again", func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case runCh <- &runReturns{
				ch: directionsCh,
			}:
			}

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case directionsCh <- types.DirectionNorth:
			}

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case sendCh <- nil:
			}

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case receiveCh <- &receiveReturns{
				data: []byte("message"),
				err:  nil,
			}:
			}

			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			case parseCh <- nil:
			}
		})

		t.Run("stop bot gracefully", func(t *testing.T) {
			botOperator.Stop()
		})
	})
}
