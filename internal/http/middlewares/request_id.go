package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := middleware.GetReqID(ctx)
		if requestID != "" {
			ctx = utils.WithField(ctx, "req_id", requestID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return middleware.RequestID(http.HandlerFunc(fn))
}
