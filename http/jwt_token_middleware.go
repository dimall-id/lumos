package http

import (
	"context"
	"net/http"
)

func JwtTokenMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			claims, err := GetTokenClaim(r.Header.Get("Authorization"))
			if err != nil {
				next.ServeHTTP(w, r)
			} else {
				ctx := context.WithValue(r.Context(), "Auth", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}


