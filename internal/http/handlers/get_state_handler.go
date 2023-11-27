package handlers

import (
	"context"
	"net/http"

	"github.com/ivan1993spb/snake-bot/internal/models"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

//counterfeiter:generate . AppGetState
type AppGetState interface {
	GetState(ctx context.Context) map[int]int
}

type GetStateHandler struct {
	app AppGetState
}

func NewGetStateHandler(app AppGetState) http.Handler {
	return &GetStateHandler{
		app: app,
	}
}

func (h *GetStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = utils.WithModule(ctx, "get_state_handler")
	log := utils.GetLogger(ctx)

	log.Info("get state handler started")

	state := h.app.GetState(ctx)
	data := models.NewGames(state)

	respond(w, r, http.StatusOK, data)
}
