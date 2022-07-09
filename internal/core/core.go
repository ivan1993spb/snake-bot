package core

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ivan1993spb/snake-bot/internal/connect"
)

var ErrRequestedTooManyBots = errors.New("requested too many bots")

type updateRequest struct {
	update update
	chres  chan<- *updateResponse
}

type updateResponse struct {
	state map[uint]uint
	err   error
}

type Connector interface {
	Connect(ctx context.Context, gameId int) (connect.Connection, error)
}

type Core struct {
	ctx context.Context

	connector Connector
	chUpdate  chan *updateRequest

	mux  *sync.Mutex
	bots map[int][]*BotOperator

	botsLimit int
}

const chUpdateBuffer = 32

// TODO: Add BotOperator factory.

func NewCore(ctx context.Context, connector Connector, botsLimit int) *Core {
	return &Core{
		ctx: ctx,

		connector: connector,
		chUpdate:  make(chan *updateRequest, chUpdateBuffer),

		mux:  &sync.Mutex{},
		bots: make(map[int][]*BotOperator),

		botsLimit: botsLimit,
	}
}

func (c *Core) getState() map[uint]uint {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.unsafeGetState()
}

func (c *Core) unsafeGetState() map[uint]uint {
	state := make(map[uint]uint, len(c.bots))
	for gameId, bots := range c.bots {
		state[gameId] = len(bots)
	}
	return state
}

const chResponseBuffer = 1

const updateTimeout = time.Millisecond * 10

func (c *Core) update(ctx context.Context, u update,
) (map[uint]uint, error) {
	chres := make(chan *updateResponse, chResponseBuffer)
	defer close(chres)

	request := &updateRequest{
		update: u,
		chres:  chres,
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	select {
	case c.chUpdate <- request:
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "update send")
	}

	var (
		response *updateResponse
		ok       bool
	)

	select {
	case response, ok = <-chres:
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "update receive")
	}

	if !ok {
		return nil, errors.New("result chan closed")
	}

	if response.err != nil {
		return nil, response.err
	}

	return response.state, nil
}

const sendResultTimeout = time.Millisecond

func (c *Core) ListenAndExecute(ctx context.Context) {
	go func() {
		for {
			var (
				request *updateRequest
				ok      bool
			)
			select {
			case request, ok = <-c.chUpdate:
			case <-ctx.Done():
				return
			}

			if !ok {
				return
			}

			state, err := c.applyUpdate(request.update)
			response := &updateResponse{
				state: state,
				err:   err,
			}
			t := time.NewTimer(sendResultTimeout)
			select {
			case request.chres <- response:
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
			t.Stop()
		}
	}()
}

func (c *Core) applyUpdate(u update) (map[uint]uint, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.unsafeApplyUpdate(u)
}

func (c *Core) unsafeApplyUpdate(u update) (map[uint]uint, error) {
	d := u.diff(c.unsafeGetState())
	for gameId, bots := range d {
		if bots > 0 {
			for i := 0; i < bots; i++ {
				bo := NewBotOperator(c.ctx, gameId, c.connector)
				bo.Run(c.ctx)
				c.bots[gameId] = append(c.bots[gameId], bo)
			}
		} else {
			for i := 0; i < -bots; i++ {
				var bo *BotOperator
				bo, c.bots[gameId] = c.bots[gameId][0], c.bots[gameId][1:]
				bo.Stop()
			}
		}
	}
	return c.unsafeGetState(), nil
}

func (c *Core) ApplyState(ctx context.Context, state State,
) (State, error) {
	return c.update(ctx, &updateBulk{
		state: state,
	})

	if stateBotsNumber(state) > c.botsLimit {
		return nil, ErrRequestedTooManyBots
	}

	//return c.unsafeGetState(), nil
	return nil, nil
}

func stateBotsNumber(state State) (number int) {
	for _, bots := range state {
		number += bots
	}
	return
}

func (c *Core) SetupOne(ctx context.Context, gameId, bots int,
) (State, error) {
	return c.update(ctx, &updateOne{
		game: gameId,
		bots: bots,
	})

	state := c.unsafeGetState()
	state[gameId] = bots

	if stateBotsNumber(state) > c.botsLimit {
		return nil, ErrRequestedTooManyBots
	}

	return c.unsafeGetState(), nil
}

func (c *Core) ReadState() State {
	return c.getState()
}
