package handlers

import (
	"net/http"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

var welcomeMessage = []byte(`Welcome to Snake-Bot!`)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = utils.WithModule(ctx, "welcome_handler")
	log := utils.GetLogger(ctx)

	log.Info("welcome handler")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(welcomeMessage); err != nil {
		log.WithError(err).Error("welcome handler fail")
	}
}
