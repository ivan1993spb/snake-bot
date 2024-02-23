package core

import (
	"context"
	"sync"
	"time"

	"github.com/ivan1993spb/snake-bot/internal/bot"
	"github.com/ivan1993spb/snake-bot/internal/connect"
	"github.com/ivan1993spb/snake-bot/internal/parser"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type BotOperator struct {
	gameId int

	connector Connector

	game   *bot.Game
	bot    *bot.Bot
	parser *parser.Parser

	stop chan struct{}
	once sync.Once
}

func NewBotOperator(ctx context.Context, gameId int,
	connector Connector) *BotOperator {
	g := bot.NewGame()
	b := bot.NewDijkstrasBot(g)
	p := &parser.Parser{
		Countdown: b,
		Me:        b,
		Size:      g,
		Game:      g,
		Printer:   utils.NewPrinterLogger(utils.GetLogger(ctx)),
	}

	return &BotOperator{
		gameId: gameId,

		connector: connector,

		game:   g,
		bot:    b,
		parser: p,

		stop: make(chan struct{}),
		once: sync.Once{},
	}
}

const botConnectRetryTimeout = time.Second * 10

func (bo *BotOperator) Run(ctx context.Context) {
	go func() {

	retry:
		conn, err := bo.connector.Connect(ctx, bo.gameId)
		if err != nil {
			utils.GetLogger(ctx).WithError(err).Error("connect fail")

			select {
			case <-bo.stop:
				return
			case <-ctx.Done():
				return
			case <-time.After(botConnectRetryTimeout):
			}

			goto retry
		}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for direction := range bo.bot.Run(ctx) {
				message := direction.ToMessageSnakeCommand()
				if err := conn.Send(ctx, message); err != nil {
					utils.GetLogger(ctx).WithError(err).Error(
						"connection send fail")
				}
			}
		}()

		go func() {
			defer wg.Done()
			for {
				message, err := conn.Receive(ctx)
				if err == connect.ErrWrongMsgType {
					// With a wrong message type we just skip the
					// message.
					continue
				}
				if err != nil {
					if err != connect.ErrConnectionClosed {
						utils.GetLogger(ctx).WithError(err).Error(
							"connection receive fail")
					}
					break
				}
				if err := bo.parser.Parse(message); err != nil {
					utils.GetLogger(ctx).WithError(err).Error(
						"parse message fail")
				}
			}
		}()

		select {
		case <-bo.stop:
		case <-ctx.Done():
		}

		if err := conn.Close(ctx); err != nil {
			utils.GetLogger(ctx).WithError(err).Error("connection close fail")
		}

		wg.Wait()
		bo.Stop()
	}()
}

func (bo *BotOperator) Stop() {
	bo.once.Do(func() {
		close(bo.stop)
	})
}
