package core

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/ivan1993spb/snake-bot/internal/connect"
)

type Connector interface {
	Connect(ctx context.Context, gameId int) (connect.Connection, error)
}

type Core struct {
	ctx context.Context

	connector Connector

	mux  *sync.Mutex
	bots map[int][]*BotOperator

	botsLimit int
}

func NewCore(ctx context.Context, connector Connector, botsLimit int) *Core {
	return &Core{
		ctx: ctx,

		connector: connector,

		mux:  &sync.Mutex{},
		bots: make(map[int][]*BotOperator),

		botsLimit: botsLimit,
	}
}

func (c *Core) state() map[int]int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.unsafeState()
}

func (c *Core) unsafeState() map[int]int {
	state := make(map[int]int, len(c.bots))
	for gameId, bots := range c.bots {
		state[gameId] = len(bots)
	}
	return state
}

var errRequestedTooManyBots = errors.New("requested too many bots")

func (c *Core) ApplyState(state map[int]int) (map[int]int, error) {
	if stateBotsNumber(state) > c.botsLimit {
		return nil, errRequestedTooManyBots
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	d := diff(c.unsafeState(), state)
	c.unsafeApplyDiff(d)

	return c.unsafeState(), nil
}

func (c *Core) unsafeApplyDiff(d map[int]int) {
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
}

func (c *Core) SetupOne(gameId, bots int) (map[int]int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	state := c.unsafeState()
	state[gameId] = bots

	if stateBotsNumber(state) > c.botsLimit {
		return nil, errRequestedTooManyBots
	}

	d := diff(c.unsafeState(), state)
	c.unsafeApplyDiff(d)

	return c.unsafeState(), nil
}

func (c *Core) ReadState() map[int]int {
	return c.state()
}
