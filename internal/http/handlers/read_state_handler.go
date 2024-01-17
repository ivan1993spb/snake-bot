package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/ivan1993spb/snake-bot/internal/models"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type AppReadState interface {
	ReadState() map[int]int
}

func NewReadStateHandler(app AppReadState) http.HandlerFunc {
	const (
		typeJson = "application/json"
		typeYaml = "text/yaml"
	)
	return func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		state := app.ReadState()
		response := models.NewGames(state)

		if accept == typeYaml {
			w.Header().Set("Content-Type", typeYaml)
			w.WriteHeader(http.StatusOK)

			err := yaml.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Log(r.Context()).WithError(err).Error("write yaml response")
			}
			return
		}

		w.Header().Set("Content-Type", typeJson)
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Log(r.Context()).WithError(err).Error("write json response")
		}
	}
}
