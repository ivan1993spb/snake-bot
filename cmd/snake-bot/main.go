package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/afero"

	"github.com/ivan1993spb/snake-bot/internal/app"
	"github.com/ivan1993spb/snake-bot/internal/config"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.StdConfig()

	ctx = utils.WithLogger(ctx, utils.NewLogger(cfg.Log))
	ctx = utils.WithModule(ctx, "main")
	log := utils.GetLogger(ctx)

	if err != nil {
		log.WithError(err).Fatal("config fail")
	}

	application := &app.App{
		Config: cfg,
		Fs:     afero.NewOsFs(),
		Clock:  utils.RealClock,
		Rand:   utils.NewRand(utils.RealClock),
	}

	application.Run(ctx)
}
