package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-bot/internal/config"
	"github.com/ivan1993spb/snake-bot/internal/connect"
	"github.com/ivan1993spb/snake-bot/internal/core"
	"github.com/ivan1993spb/snake-bot/internal/secure"
	"github.com/ivan1993spb/snake-bot/internal/server"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

const ApplicationName = "Snake-Bot"

var (
	Version = "dev"
	Build   = "dev"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.StdConfig()

	{
		logger := utils.NewLogger(cfg.Log)
		ctx = utils.LogContext(ctx, logger)
	}

	if err != nil {
		utils.Log(ctx).WithError(err).Fatal("config fail")
	}

	utils.Log(ctx).WithFields(logrus.Fields{
		"version": Version,
		"build":   Build,
	}).Info("Welcome to Snake-Bot!")

	sec := secure.NewSecure()
	if err := sec.GenerateToken(os.Stdout); err != nil {
		utils.Log(ctx).WithError(err).Fatal("security fail")
	}
	utils.Log(ctx).Warn("auth token successfully generated")

	headerAppInfo := utils.FormatAppInfoHeader(ApplicationName, Version, Build)

	connector := connect.NewConnector(cfg.Target, headerAppInfo)

	c := core.NewCore(ctx, connector, cfg.Bots.Limit)

	serv := server.NewServer(ctx, cfg.Server, headerAppInfo, c, sec)

	if err := serv.ListenAndServe(ctx); err != nil {
		utils.Log(ctx).WithError(err).Fatal("server error")
	}

	utils.Log(ctx).Info("buh bye!")
}
