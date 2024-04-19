package handlers

import (
	"net/http"

	snakebot "github.com/ivan1993spb/snake-bot"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

func OpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = utils.WithModule(ctx, "openapi_handler")
	log := utils.GetLogger(ctx)

	log.Info("openapi handler")

	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(snakebot.OpenAPISpec)); err != nil {
		log.WithError(err).Error("openapi handler fail")
	}
}
