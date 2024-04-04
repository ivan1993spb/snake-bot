package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5/request"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type Secure interface {
	VerifyToken(tokenString string) (string, error)
}

func JwtTokenAuth(sec Secure) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := utils.GetLogger(ctx)

			tokenString, err := request.BearerExtractor{}.ExtractToken(r)
			if err != nil {
				log.WithError(err).Error("error extracting token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			subject, err := sec.VerifyToken(tokenString)
			if err != nil {
				log.WithError(err).Error("error verifying token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = utils.WithField(ctx, "subject", subject)
			log = utils.GetLogger(ctx)

			log.Info("token verified")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
