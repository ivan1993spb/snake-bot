package app

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/ivan1993spb/snake-bot/internal/config"
	"github.com/ivan1993spb/snake-bot/internal/connect"
	"github.com/ivan1993spb/snake-bot/internal/core"
	"github.com/ivan1993spb/snake-bot/internal/http"
	"github.com/ivan1993spb/snake-bot/internal/secure"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type App struct {
	Config config.Config
	Fs     afero.Fs
	Clock  utils.Clock
	Rand   core.Rand
}

const shutdownTimeout = time.Second * 5

func (a *App) Run(ctx context.Context) {
	log := utils.GetLogger(ctx)

	log.WithFields(logrus.Fields{
		"version": Version,
		"build":   Build,
	}).Info("Welcome to Snake-Bot!")

	// This is a header that will contain information about version and build.
	headerAppInfo := utils.FormatAppInfoHeader(ApplicationName, Version, Build)

	// Module "secure" is responsible for authentication.
	jwtSec, err := secure.New(a.Fs, a.Clock).JwtFromFile(ctx, a.Config.Server.JWTSecret)
	if err != nil {
		log.WithError(err).Fatal("security fail")
	}

	// Module "connect" is responsible for connecting to the target server.
	connector := connect.NewConnector(a.Config.Target, headerAppInfo)

	// factory creates bot operators which are responsible for
	// managing bots and their sessions.
	factory := &core.DijkstrasBotOperatorFactory{
		Logger:    utils.GetLogger(utils.WithModule(ctx, "notification")),
		Rand:      a.Rand,
		Connector: connector,
		Clock:     a.Clock,
	}

	// Module "core" manages bot operators.
	appCore := core.NewCore(&core.Params{
		BotsLimit:          a.Config.Bots.Limit,
		BotOperatorFactory: factory,
		Clock:              a.Clock,
	})

	done := appCore.Run(utils.WithModule(ctx, "core"))

	// Start the REST API server.
	server := http.NewServer(http.ServerParams{
		Config:  a.Config.Server,
		AppInfo: headerAppInfo,
		Core:    appCore,
		Secure:  jwtSec,
	})

	err = server.ListenAndServe(utils.WithModule(ctx, "server"))
	if err != nil {
		log.WithError(err).Fatal("server fail")
	}

	select {
	case <-done:
	case <-time.After(shutdownTimeout):
		log.Fatal("kill")
	}

	log.Info("buh bye!")
}
