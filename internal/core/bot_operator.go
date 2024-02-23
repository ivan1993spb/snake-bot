package core

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/ivan1993spb/snake-bot/internal/bot"
	"github.com/ivan1993spb/snake-bot/internal/connect"
	"github.com/ivan1993spb/snake-bot/internal/parser"
	"github.com/ivan1993spb/snake-bot/internal/types"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

//counterfeiter:generate . BotEngine
type BotEngine interface {
	Run(ctx context.Context) <-chan types.Direction
}

//counterfeiter:generate . Connector
type Connector interface {
	Connect(ctx context.Context, gameId int) (connect.Connection, error)
}

//counterfeiter:generate . Rand
type Rand interface {
	Intn(n int) int
}

//counterfeiter:generate . Parser
type Parser interface {
	Parse(message []byte) error
}

var _ BotOperator = (*botOperator)(nil)

type botOperator struct {
	gameId int

	connector Connector
	bot       BotEngine
	parser    Parser
	rand      Rand
	clock     utils.Clock

	stop chan struct{}
	once sync.Once
}

type BotOperatorParams struct {
	GameId    int
	Connector Connector
	BotEngine BotEngine
	Parser    Parser
	Rand      Rand
	Clock     utils.Clock
}

type DijkstrasBotOperatorFactory struct {
	Logger    *logrus.Entry
	Rand      Rand
	Connector Connector
	Clock     utils.Clock
}

func (f *DijkstrasBotOperatorFactory) New(gameId int) BotOperator {
	g := bot.NewGame()
	b := bot.NewDijkstrasBot(g)
	p := &parser.Parser{
		Countdown: b,
		Me:        b,
		Size:      g,
		Game:      g,
		Printer:   utils.NewPrinterLogger(f.Logger.WithField("game", gameId)),
	}

	return NewBotOperator(&BotOperatorParams{
		GameId:    gameId,
		Connector: f.Connector,
		BotEngine: b,
		Parser:    p,
		Rand:      f.Rand,
		Clock:     f.Clock,
	})
}

func NewBotOperator(params *BotOperatorParams) BotOperator {
	return &botOperator{
		gameId: params.GameId,

		connector: params.Connector,
		bot:       params.BotEngine,
		parser:    params.Parser,
		rand:      params.Rand,
		clock:     params.Clock,

		stop: make(chan struct{}),
		once: sync.Once{},
	}
}

func (bo *botOperator) Run(ctx context.Context) {
	// A context for the whole bot operator.
	ctx = utils.WithModule(ctx, "operator")
	ctx = utils.WithTaskId(ctx)
	ctx, cancel := context.WithCancel(ctx)

	log := utils.GetLogger(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		select {
		case <-bo.stop:
			// Stopping the current bot operator.
			cancel()
		case <-ctx.Done():
			// When the whole applictaion is shutting down.
			return
		}
	}()

	go func() {
		defer wg.Done()

		log.Info("bot operator started")
		defer log.Info("bot operator stopped")

		err := bo.runSession(ctx)
		if err != nil {
			log.WithError(err).Error("session failure")
		}
	}()

	wg.Wait()

	// Call the stop function anyway.
	bo.Stop()
}

const botConnectRetryTimeout = time.Second * 15

// tryConnect tries to connect to the server. It retries until the context is
// done.
func (bo *botOperator) tryConnect(ctx context.Context) (connect.Connection, error) {
	log := utils.GetLogger(ctx)

	for {
		conn, err := bo.connector.Connect(ctx, bo.gameId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil, err
			}

			log.WithError(err).Error("connecting failure")

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-bo.clock.After(botConnectRetryTimeout):
			}

			log.Info("retrying connection")
			continue
		}

		log.Info("connected")

		return conn, nil
	}
}

const (
	botMinConnectionDelayMs = 3000
	botMaxConnectionDelayMs = 15000
)

func (bo *botOperator) randomDelay(ctx context.Context) error {
	log := utils.GetLogger(ctx)

	ms := bo.rand.Intn(botMaxConnectionDelayMs - botMinConnectionDelayMs)
	ms += botMinConnectionDelayMs

	delay := time.Duration(ms) * time.Millisecond

	log.WithField("delay", delay).Info("delaying connection")

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-bo.clock.After(delay):
	}

	return nil
}

// runSession runs a session of the bot operator. It recreates bots and
// reconnects as needed.
func (bo *botOperator) runSession(ctx context.Context) error {
	ctx = utils.WithModule(ctx, "session")
	log := utils.GetLogger(ctx)

	log.Info("starting session")
	defer log.Info("session stopped")

	for {
		if err := bo.randomDelay(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		}

		conn, err := bo.tryConnect(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		}

		botErr := bo.runBot(ctx, conn)

		log.Info("closing connection")
		if err := conn.Close(ctx); err != nil {
			log.WithError(err).Error("connection close fail")
		}

		// If the context is canceled, we just return.
		if botErr == nil || errors.Is(botErr, context.Canceled) {
			return nil
		}

		log.Info("restarting bot")
	}
}

// runBot runs 1 bot instance within a session.
func (bo *botOperator) runBot(ctx context.Context, conn connect.Connection) error {
	ctx, cancel := context.WithCancel(ctx)
	ctx = utils.WithModule(ctx, "bot")
	log := utils.GetLogger(ctx)

	log.Info("starting bot")
	defer log.Info("bot stopped")

	g := new(errgroup.Group)

	g.Go(func() error {
		defer cancel()

		defer log.Debug("sender stopped")

		for direction := range bo.bot.Run(ctx) {
			message := direction.ToMessageSnakeCommand()
			if err := conn.Send(ctx, message); err != nil {
				log.WithError(err).Error("failed to send")
				if errors.Is(err, context.Canceled) {
					return nil
				}

				return err
			}
		}

		return nil
	})

	g.Go(func() error {
		defer cancel()

		defer log.Debug("receiver stopped")

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			message, err := conn.Receive(ctx)
			if err == connect.ErrWrongMsgType {
				// With a wrong message type we just skip the
				// message.
				continue
			}

			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				if err != connect.ErrConnectionClosed {
					log.WithError(err).Error("receive fail")
					return err
				}
				break
			}

			if err := bo.parser.Parse(message); err != nil {
				log.WithError(err).Error("parse fail")
			}
		}

		return nil
	})

	return g.Wait()
}

func (bo *botOperator) Stop() {
	bo.once.Do(func() {
		close(bo.stop)
	})
}
