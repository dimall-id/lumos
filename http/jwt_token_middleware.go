package http

import (
	"context"
	"github.com/dimall-id/lumos/v2/misc"
	"net/http"
)

func GetTokenClaim (authentication string) (AccessToken, error) {
	tokens := misc.BuildToMap(`Bearer (?P<token>[\W\w]+)`, authentication)
	t := AccessToken{}
	err := t.FromJwtBase64(tokens["token"])
	if err != nil {return AccessToken{}, err}
	return t, nil
}

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


