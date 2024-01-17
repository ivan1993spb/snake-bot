package middlewares

import "net/http"

const headerToken = "X-Snake-Bot-Token"

type Secure interface {
	VerifyToken(token string) bool
}

func TokenAuth(sec Secure) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token := r.Header.Get(headerToken); !sec.VerifyToken(token) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
