package http

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ReqIdMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New().String()
		r.WithContext(context.WithValue(r.Context(), "Req-Id", reqId))
		log.WithField("Req-Id", reqId)
		next.ServeHTTP(w, r)
	})
}