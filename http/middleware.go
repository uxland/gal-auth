package http

import (
	"github.com/uxland/ga-go-auth/shared"
	"net/http"
	"strings"
)

func MiddlewareFactory(apiSecret string, scheme shared.AuthenticationScheme) func(handler http.Handler) http.Handler {
	secret := []byte(apiSecret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
			}
			authHeader := r.Header.Get(shared.AuthorizationHeader)
			if authHeader == "" {
				next.ServeHTTP(w, r)
			}
			tokenSplit := strings.Split(authHeader, " ")
			bearer := tokenSplit[1]
			token, err := shared.VerifyToken(bearer, secret)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			ctx := shared.SetClaimsForContext(r.Context(), token)
			r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
