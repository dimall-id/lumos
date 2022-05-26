package http

import (
	"net/http"
	"testing"
)

func TestAddMiddleware(t *testing.T) {
	AddMiddleware(func(handler http.Handler) http.Handler {
		return handler
	})

	if len(middlewares) != 1 {
		t.Error("Fail to test, middleware doesn't added to middlewares")
	}
}


func TestGetAllMiddleware(t *testing.T) {
	mws := GetAllMiddleware()

	if len(mws) != 1 {
		t.Error("Fail to test, middleware doesn't return right len")
	}
}