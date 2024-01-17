package handlers

import (
	"net/http"

	"github.com/ivan1993spb/snake-bot"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

func OpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(snakebot.OpenAPISpec)); err != nil {
		utils.Log(r.Context()).WithError(err).Error("openapi handler fail")
	}
}
