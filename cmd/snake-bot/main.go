package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-bot/internal/config"
	"github.com/ivan1993spb/snake-bot/internal/connect"
	"github.com/ivan1993spb/snake-bot/internal/core"
	"github.com/ivan1993spb/snake-bot/internal/http"
	"github.com/ivan1993spb/snake-bot/internal/secure"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

const ApplicationName = "Snake-Bot"

var (
	Version = "dev"
	Build   = "dev"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.StdConfig()

	{
		logger := utils.NewLogger(cfg.Log)
		ctx = utils.WithLogger(ctx, logger)
	}

	if err != nil {
		utils.GetLogger(ctx).WithError(err).Fatal("config fail")
	}

	utils.GetLogger(ctx).WithFields(logrus.Fields{
		"version": Version,
		"build":   Build,
	}).Info("Welcome to Snake-Bot!")

	sec := secure.NewSecure()
	if err := sec.GenerateToken(os.Stdout); err != nil {
		utils.GetLogger(ctx).WithError(err).Fatal("security fail")
	}
	utils.GetLogger(ctx).Warn("auth token successfully generated")

	headerAppInfo := utils.FormatAppInfoHeader(ApplicationName, Version, Build)

	// Module "connect" is responsible for connecting to the target server.
	connector := connect.NewConnector(cfg.Target, headerAppInfo)

	rand := utils.NewRand(utils.RealClock)

	// factory creates bot operators which are responsible for
	// managing bots and their sessions.
	factory := &core.DijkstrasBotOperatorFactory{
		Logger:    utils.GetLogger(utils.WithModule(ctx, "notification")),
		Rand:      rand,
		Connector: connector,
		Clock:     utils.RealClock,
	}

	// Module "core" manages bot operators.
	c := core.NewCore(&core.Params{
		BotsLimit:          cfg.Bots.Limit,
		BotOperatorFactory: factory,
		Clock:              utils.RealClock,
	})

	serv := http.NewServer(ctx, cfg.Server, headerAppInfo, c, sec)

	if err := serv.ListenAndServe(ctx); err != nil {
		utils.GetLogger(ctx).WithError(err).Fatal("server error")
	}

	utils.GetLogger(ctx).Info("buh bye!")
}
