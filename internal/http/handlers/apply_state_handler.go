package handlers

import (
	"context"
	"encoding/json"
	"mime"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/ivan1993spb/snake-bot/internal/models"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type AppApplyState interface {
	SetState(ctx context.Context, state map[int]int) (map[int]int, error)
	SetOne(ctx context.Context, gameId, botsNumber int) (map[int]int, error)
}

func NewApplyStateHandler(app AppApplyState) http.HandlerFunc {
	const (
		mediaTypeFormUrlencoded = "application/x-www-form-urlencoded"
		mediaTypeJson           = "application/json"
		mediaTypeYaml           = "text/yaml"
	)

	applyState := func(games *models.Games) (map[int]int, error, int) {
		requested := games.ToMapState()
		state, err := app.SetState(context.TODO(), requested)
		if err != nil {
			return nil, errors.Wrap(err, "apply state fail"),
				http.StatusInternalServerError
		}
		return state, nil, http.StatusCreated
	}

	setupOne := func(gameId, bots int) (map[int]int, error, int) {
		state, err := app.SetOne(context.TODO(), gameId, bots)
		if err != nil {
			return nil, errors.Wrap(err, "setup one fail"),
				http.StatusInternalServerError
		}
		return state, err, http.StatusCreated
	}

	handleFormUrlencoded := func(r *http.Request) (map[int]int, error, int) {
		if err := r.ParseForm(); err != nil {
			return nil, errors.Wrap(err, "parse form fail"),
				http.StatusBadRequest
		}
		gameId, err := strconv.Atoi(r.PostForm.Get("game"))
		if err != nil {
			return nil, errors.Wrap(err, "parse game id"),
				http.StatusBadRequest
		}
		bots, err := strconv.Atoi(r.PostForm.Get("bots"))
		if err != nil {
			return nil, errors.Wrap(err, "parse bots number"),
				http.StatusBadRequest
		}
		return setupOne(gameId, bots)
	}

	handleJson := func(r *http.Request) (map[int]int, error, int) {
		var games *models.Games
		if err := json.NewDecoder(r.Body).Decode(&games); err != nil {
			return nil, errors.Wrap(err, "decode json fail"),
				http.StatusBadRequest
		}
		return applyState(games)
	}

	handleYaml := func(r *http.Request) (map[int]int, error, int) {
		var games *models.Games
		if err := yaml.NewDecoder(r.Body).Decode(&games); err != nil {
			return nil, errors.Wrap(err, "decode yaml fail"),
				http.StatusBadRequest
		}
		return applyState(games)
	}

	handlersMapping := map[string]func(r *http.Request) (map[int]int, error, int){
		mediaTypeFormUrlencoded: handleFormUrlencoded,
		mediaTypeJson:           handleJson,
		mediaTypeYaml:           handleYaml,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-type")
		acceptType := r.Header.Get("Accept")
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			utils.GetLogger(r.Context()).WithError(err).Error("parse media type")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handler, ok := handlersMapping[mediaType]
		if !ok {
			utils.GetLogger(r.Context()).WithField(
				"media_type", mediaType).Error("unknown media type")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		state, err, statusCode := handler(r)
		if err != nil {
			utils.GetLogger(r.Context()).WithError(err).Error("handle fail")
			w.WriteHeader(statusCode)
			return
		}

		response := models.NewGames(state)

		if acceptType == mediaTypeYaml {
			w.Header().Set("Content-Type", mediaTypeYaml)
			w.WriteHeader(statusCode)
			if err := yaml.NewEncoder(w).Encode(response); err != nil {
				utils.GetLogger(r.Context()).WithError(err).Error("write yaml response")
			}
			return
		}

		w.Header().Set("Content-Type", mediaTypeJson)
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			utils.GetLogger(r.Context()).WithError(err).Error("write json response")
		}
	}
}
