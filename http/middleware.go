package http

import "github.com/gorilla/mux"

var middlewares []mux.MiddlewareFunc

func AddMiddleware (middleware mux.MiddlewareFunc) {
	middlewares = append(middlewares, middleware)
}

func GetAllMiddleware () []mux.MiddlewareFunc {
	return middlewares
}