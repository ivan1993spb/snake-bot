package core

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

//counterfeiter:generate . BotOperator
type BotOperator interface {
	Run(ctx context.Context)
	Stop()
}

//counterfeiter:generate . BotOperatorFactory
type BotOperatorFactory interface {
	New(gameId int) BotOperator
}

type stateRquest struct {
	state  map[int]int
	result chan<- map[int]int
}

type Core struct {
	mux  sync.Mutex
	wg   sync.WaitGroup
	bots map[int][]BotOperator

	botsLimit int

	applyStateCh chan *stateRquest

	factory BotOperatorFactory
	clock   utils.Clock

	storage Storage
}

type Params struct {
	BotsLimit int
	BotOperatorFactory
	utils.Clock
	Storage
}

const applyStateChSize = 100

func NewCore(params *Params) *Core {
	return &Core{
		bots: make(map[int][]BotOperator),

		botsLimit: params.BotsLimit,

		applyStateCh: make(chan *stateRquest, applyStateChSize),

		factory: params.BotOperatorFactory,
		clock:   params.Clock,

		storage: params.Storage,
	}
}

const sendResultTimeout = time.Millisecond * 10

func (c *Core) Run(ctx context.Context) <-chan struct{} {
	log := utils.GetLogger(ctx)

	done := make(chan struct{})

	go func() {
		// Signal the caller that the core is stopped.
		defer close(done)

		defer func() {
			close(c.applyStateCh)
			log.Info("core stopped")
		}()

		defer func() {
			// Wait the bots to stop
			c.wg.Wait()
			log.Info("all bots stopped")
		}()

		log.Info("core started")

		// TODO: Consider returning error from preloadState
		// Preload state from storage
		c.preloadState(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case req := <-c.applyStateCh:
				log.Debug("applying new state")

				// TODO: Consider returning error from applyState and
				//       sending it to the caller.
				result := c.applyState(ctx, req.state)
				c.sendResult(ctx, req.result, result)
			}
		}
	}()

	return done
}

func (c *Core) preloadState(ctx context.Context) {
	log := utils.GetLogger(ctx)

	log.Info("loading state from storage")

	state, err := c.storage.Load(ctx)
	if err != nil {
		log.WithError(err).Error("failed to load state from storage")
		return
	}

	if stateBotsNumber(state) > c.botsLimit {
		log.WithField("bots_limit", c.botsLimit).Error("loaded state exceeds bots limit")
		return
	}

	c.applyState(ctx, state)
}

func (c *Core) applyState(ctx context.Context, state map[int]int) map[int]int {
	c.mux.Lock()
	defer c.mux.Unlock()

	log := utils.GetLogger(ctx)

	log.Info("applying new state")

	oldState := c.unsafeGetState()
	// Only the diff is applied to the state to avoid unnecessary
	// restarts
	d := diff(oldState, state)
	if len(d) == 0 {
		log.Info("no changes in state")
		return state
	}

	add, remove := diffStats(d)
	log.WithFields(logrus.Fields{
		"add":    add,
		"remove": remove,
	}).Info("applying diff to the current state")
	c.unsafeApplyDiff(ctx, d)

	// Save the new state
	state = c.unsafeGetState()
	err := c.storage.Save(ctx, state)
	if err != nil {
		log.WithError(err).Error("failed to save state to storage")

		log.Info("reverting changes")
		c.unsafeApplyDiff(ctx, invertDiff(d))

		return oldState
	}

	return state
}

func (c *Core) sendResult(
	ctx context.Context,
	result chan<- map[int]int,
	state map[int]int,
) {

	defer close(result)

	select {
	case <-ctx.Done():
	case <-c.clock.After(sendResultTimeout):
	case result <- state:
	}
}

func (c *Core) unsafeGetState() map[int]int {
	state := make(map[int]int, len(c.bots))
	for gameId, bots := range c.bots {
		state[gameId] = len(bots)
	}
	return state
}

var ErrRequestedTooManyBots = errors.New("requested too many bots")

func (c *Core) SetState(ctx context.Context, state map[int]int) (map[int]int, error) {
	if stateBotsNumber(state) > c.botsLimit {
		return nil, ErrRequestedTooManyBots
	}

	ch := make(chan map[int]int, 1)

	req := &stateRquest{
		state:  state,
		result: ch,
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case c.applyStateCh <- req:
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case state := <-ch:
		return state, nil
	}
}

func (c *Core) unsafeApplyDiff(ctx context.Context, d map[int]int) {
	for gameId, bots := range d {
		if bots > 0 {
			c.unsafeSpawn(ctx, gameId, bots)
		} else {
			c.unsafeTerminate(ctx, gameId, -bots)
		}
	}
}

func (c *Core) unsafeSpawn(ctx context.Context, gameId, bots int) {
	c.wg.Add(bots)

	for i := 0; i < bots; i++ {
		// Initialize new bot
		bot := c.factory.New(gameId)

		// Start bot
		go func(gameId int) {
			defer c.wg.Done()

			bot.Run(utils.WithField(ctx, "game", gameId))
		}(gameId)

		c.bots[gameId] = append(c.bots[gameId], bot)
	}
}

func (c *Core) unsafeTerminate(ctx context.Context, gameId, bots int) {
	for i := 0; i < bots && len(c.bots[gameId]) > 0; i++ {
		c.bots[gameId][0].Stop()
		c.bots[gameId] = c.bots[gameId][1:]
	}

	if len(c.bots[gameId]) == 0 {
		delete(c.bots, gameId)
	}
}

func (c *Core) SetOne(ctx context.Context, gameId, bots int) (map[int]int, error) {
	state := c.GetState(ctx)
	state[gameId] = bots

	if stateBotsNumber(state) > c.botsLimit {
		return nil, ErrRequestedTooManyBots
	}

	return c.SetState(ctx, state)
}

func (c *Core) GetState(ctx context.Context) map[int]int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.unsafeGetState()
}
