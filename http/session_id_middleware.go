package http

import (
	log "github.com/dimall-id/lumos/v2/logger"
	"net/http"
)

func SessionIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Session-Id") != "" {
			log.AddField("Session-Id", r.Header.Get("X-Session-Id"))
		}
	})
}