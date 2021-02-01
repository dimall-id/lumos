package http

import (
	"encoding/json"
	"github.com/dimall-id/lumos/http/route"
	"github.com/gorilla/mux"
	"net/http"
)

func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		err := MethodNotAllow()
		res, _ := json.Marshal(err)
		w.Write(res)
	})
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		err := NotFound()
		res, _ := json.Marshal(err)
		w.Write(res)
	})
}

func handleRequest(w http.ResponseWriter, r *http.Request, f func(r2 *http.Request) (interface{}, HttpError)) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	data, err := f(r)
	var res []byte
	if err.Message != "" {
		w.WriteHeader(err.Code)
		res, _ = json.Marshal(err)
	} else {
		res, _ = json.Marshal(data)
	}

	w.Write(res)
}

func generateMuxRouter () *mux.Router {
	r := mux.NewRouter()
	r.MethodNotAllowedHandler = methodNotAllowedHandler()
	r.NotFoundHandler = notFoundHandler()

	for i, _ := range route.GetAll() {
		rr := route.GetAt(i)
		r.HandleFunc(rr.Url, func(w http.ResponseWriter, r *http.Request) {
			handleRequest(w, r, rr.Func)
		}).Methods(rr.HttpMethod).Name(rr.Name)
	}

	return r
}

func StartHttpServer(port string) error {
	r := generateMuxRouter()
	return http.ListenAndServe(port, r)
}