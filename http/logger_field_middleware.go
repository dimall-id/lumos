package http

import (
	"context"
	"net/http"
)


// This middleware's job is to collect header data that will be use as http request logger default field.
func LoggerFieldMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		fields := map[string]interface{}{}
		if r.Header.Get("X-Session-Id") != "" {
			fields["X-Session-Id"] = r.Header.Get("X-Session-Id")
		}
		ctx := context.WithValue(r.Context(), "LoggerFields", fields)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}