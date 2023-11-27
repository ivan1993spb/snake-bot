package handlers

import (
	"context"
	"encoding/json"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/ivan1993spb/snake-bot/internal/core"
	"github.com/ivan1993spb/snake-bot/internal/models"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

//counterfeiter:generate . AppSetState
type AppSetState interface {
	SetState(ctx context.Context, state map[int]int) (map[int]int, error)
	SetOne(ctx context.Context, gameId, botsNumber int) (map[int]int, error)
}

type SetStateHandler struct {
	app AppSetState
}

func NewSetStateHandler(app AppSetState) http.Handler {
	return &SetStateHandler{
		app: app,
	}
}

const (
	mediaTypeFormUrlencoded = "application/x-www-form-urlencoded"
	mediaTypeJson           = "application/json"
	mediaTypeYaml           = "text/yaml"
)

const setStateTimeout = 200 * time.Millisecond

func (h *SetStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = utils.WithModule(ctx, "set_state_handler")
	log := utils.GetLogger(ctx)

	log.Info("set state handler started")

	ctx, cancel := context.WithTimeout(ctx, setStateTimeout)
	defer cancel()

	contentType := r.Header.Get("Content-type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.WithError(err).Error("parse media type")

		respondError(w, r, http.StatusBadRequest)
		return
	}

	log = log.WithField("media_type", mediaType)
	ctx = utils.WithLogger(ctx, log)
	r = r.WithContext(ctx)

	var (
		state map[int]int
		code  int
	)

	log.Info("process request")

	switch mediaType {
	case mediaTypeFormUrlencoded:
		state, code, err = h.handleFormUrlencoded(w, r)
	case mediaTypeJson:
		state, code, err = h.handleJson(w, r)
	case mediaTypeYaml:
		state, code, err = h.handleYaml(w, r)
	default:
		log.Error("invalid media type")

		respondError(w, r, http.StatusUnsupportedMediaType)
		return
	}

	if err != nil {
		log.WithError(err).Error("handle request")

		respondError(w, r, code)
		return
	}

	data := models.NewGames(state)

	respond(w, r, http.StatusCreated, data)
}

func appSetStateErrStatus(err error) int {
	if errors.Is(err, core.ErrRequestedTooManyBots) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func (h *SetStateHandler) handleFormUrlencoded(
	w http.ResponseWriter,
	r *http.Request,
) (
	map[int]int,
	int,
	error,
) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, "parse form fail")
	}

	gameId, err := strconv.Atoi(r.PostForm.Get("game"))
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, "parse game id fail")
	}

	bots, err := strconv.Atoi(r.PostForm.Get("bots"))
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, "parse bots number fail")
	}

	state, err := h.app.SetOne(ctx, gameId, bots)
	if err != nil {
		return nil, appSetStateErrStatus(err), errors.Wrap(err, "set one fail")
	}

	return state, http.StatusCreated, nil
}

func (h *SetStateHandler) handleJson(
	w http.ResponseWriter,
	r *http.Request,
) (
	map[int]int,
	int,
	error,
) {
	ctx := r.Context()

	var games *models.Games

	err := json.NewDecoder(r.Body).Decode(&games)
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, "decode json fail")
	}

	state, err := h.app.SetState(ctx, games.ToMapState())
	if err != nil {
		return nil, appSetStateErrStatus(err), errors.Wrap(err, "set state fail")
	}

	return state, http.StatusCreated, nil
}

func (h *SetStateHandler) handleYaml(
	w http.ResponseWriter,
	r *http.Request,
) (
	map[int]int,
	int,
	error,
) {
	ctx := r.Context()

	var games *models.Games

	err := yaml.NewDecoder(r.Body).Decode(&games)
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, "decode yaml fail")
	}

	state, err := h.app.SetState(ctx, games.ToMapState())
	if err != nil {
		return nil, appSetStateErrStatus(err), errors.Wrap(err, "set state fail")
	}

	return state, http.StatusCreated, nil
}
