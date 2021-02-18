package http

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := MethodNotAllow()
		w.WriteHeader(err.Code)
		res, _ := json.Marshal(err)
		var dest bytes.Buffer
		json.Compact(&dest, res)
		w.Write(dest.Bytes())
	})
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := NotFound()
		w.WriteHeader(err.Code)
		res, _ := json.Marshal(err)
		var dest bytes.Buffer
		json.Compact(&dest, res)
		w.Write(dest.Bytes())
	})
}

func HandleRequest(w http.ResponseWriter, r *http.Request, f func(r2 *http.Request) (interface{}, HttpError)) {
	w.Header().Set("Content-Type", "application/json")
	data, err := f(r)
	var res []byte
	var dest bytes.Buffer
	if err.Message != "" {
		w.WriteHeader(err.Code)
		res, _ = json.Marshal(err)
		json.Compact(&dest, res)
	} else {
		res, _ = json.Marshal(data)
		json.Compact(&dest, res)
	}

	w.Write(dest.Bytes())
}

func GenerateMuxRouter (routes []Route, middleware []mux.MiddlewareFunc) *mux.Router {
	r := mux.NewRouter()
	r.MethodNotAllowedHandler = methodNotAllowedHandler()
	r.NotFoundHandler = notFoundHandler()

	for i, _ := range routes {
		rr := GetRouteAt(i)
		r.HandleFunc(rr.Url, func(w http.ResponseWriter, r *http.Request) {
			HandleRequest(w, r, rr.Func)
		}).Methods(rr.HttpMethod).Name(rr.Name)
	}

	r.Use(ContentTypeMiddleware)
	for _, mwr := range middleware {
		mw := mwr
		r.Use(mw)
	}

	return r
}

func StartHttpServer(port string) error {
	r := GenerateMuxRouter(routes, middlewares)
	return http.ListenAndServe(port, r)
}