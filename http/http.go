package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dimall-id/lumos/misc"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/dimall-id/jwt-go"
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

func CheckRole (roles []string, routes []string) bool {
	route := make(map[string]string)
	for _,d := range routes {
		route[d] = d
	}
	for _, role := range roles {
		if _, oke := route[role]; oke {
			return true
		}
	}
	return false
}

func CheckAuthentication (authentication string, rr Route) HttpError {
	if len(rr.Roles) <= 0 {
		return HttpError{}
	} else {
		if authentication == "" {
			return Unauthorized()
		} else {
			t := misc.BuildToMap(`Bearer (?P<token>[\W\w]+)`, authentication)
			token, err := jwt.ParseUnverified(t["token"], jwt.MapClaims{})
			if err != nil {
				fmt.Println(err)
				return BadRequest()
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			fmt.Println(claims)
			if claim, oke := claims["Roles"] ; oke {
				if !CheckRole(claim.([]string), rr.Roles) {
					return Unauthorized()
				}
			} else {
				return BadRequest()
			}
		}
	}
	return HttpError{}
}

func BuildJsonResponse (response interface{}) []byte {
	res, _ := json.Marshal(response)
	var dest bytes.Buffer
	json.Compact(&dest, res)
	return dest.Bytes()
}

func HandleRequest(w http.ResponseWriter, r *http.Request, rr Route) {
	var res []byte
	err := CheckAuthentication(r.Header.Get("Authentication"), rr)
	if err.Message != "" {
		res = BuildJsonResponse(err)
	} else {
		data, err := rr.Func(r)
		if err.Message != "" {
			res = BuildJsonResponse(err)
		} else {
			res = BuildJsonResponse(data)
		}
	}
	w.Write(res)
}

func GenerateMuxRouter (routes []Route, middleware []mux.MiddlewareFunc) *mux.Router {
	r := mux.NewRouter()
	r.MethodNotAllowedHandler = methodNotAllowedHandler()
	r.NotFoundHandler = notFoundHandler()

	for i, _ := range routes {
		rr := GetRouteAt(i)
		r.HandleFunc(rr.Url, func(w http.ResponseWriter, r *http.Request) {
			HandleRequest(w, r, rr)
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