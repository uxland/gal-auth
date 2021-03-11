package http

import (
	"github.com/uxland/gal-auth/shared"
	"net/http"
	"strings"
)

func MiddlewareFactory(apiSecret string) func(handler http.Handler) http.Handler {
	secret := []byte(apiSecret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}
			authHeader := r.Header.Get(shared.AuthorizationHeader)
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenSplit := strings.Split(authHeader, " ")

			scheme := tokenSplit[0]
			auth := tokenSplit[1]

			ctx, err := shared.AuthenticateInputToken(r.Context(), scheme, auth, secret)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Unauthorized"))
				return
			}
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
