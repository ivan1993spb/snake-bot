package handlers

import (
	"net/http"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

var welcomeMessage = []byte(`Welcome to Snake-Bot!`)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(welcomeMessage); err != nil {
		utils.GetLogger(r.Context()).WithError(err).Error("welcome handler fail")
	}
}
