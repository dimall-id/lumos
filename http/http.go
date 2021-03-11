package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dimall-id/lumos/misc"
	"github.com/gorilla/mux"
	"net/http"
)

func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := MethodNotAllow()
		w.WriteHeader(err.Code)
		res, _ := json.Marshal(err)
		var dest bytes.Buffer
		json.Compact(&dest, res)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Write(dest.Bytes())
	})
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := NotFound()
		w.WriteHeader(err.Code)
		res, _ := json.Marshal(err)
		var dest bytes.Buffer
		json.Compact(&dest, res)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Write(dest.Bytes())
	})
}

/**
@TODO : Make unit testing
 */
func CheckRole (roles []string, routes map[string]string) bool {
	for _, role := range roles {
		if _, oke := routes[fmt.Sprintf("%v", role)]; oke {return true}
	}
	return false
}

/**
@TODO : make unit testing
 */
func CheckAuthorization(authentication string, rr Route) HttpError {
	roles := make(map[string]string)
	for _, role := range rr.Roles {roles[role] = role}
	if len(roles) <= 0 {return HttpError{}}
	if _, oke := roles["ANONYMOUS"]; oke {return HttpError{}}

	if authentication == "" {
		return Unauthorized()
	} else {
		claims, err := GetTokenClaim(authentication)
		if err != nil {return BadRequest()}
		if !CheckRole(claims.Roles, roles) {return Unauthorized()}
	}
	return HttpError{}
}

/**
@TODO : make unit testing
 */
func GetTokenClaim (authentication string) (AccessToken, error) {
	tokens := misc.BuildToMap(`Bearer (?P<token>[\W\w]+)`, authentication)
	t := AccessToken{}
	err := t.FromJwtBase64(tokens["token"])
	if err != nil {return AccessToken{}, err}
	return t, nil
}

func BuildJsonResponse (response interface{}) []byte {
	res, _ := json.Marshal(response)
	var dest bytes.Buffer
	json.Compact(&dest, res)
	return dest.Bytes()
}

func HandleRequest(w http.ResponseWriter, r *http.Request, rr Route) {
	var res []byte
	err := CheckAuthorization(r.Header.Get("Authorization"), rr)
	if err.Message != "" {
		res = BuildJsonResponse(err)
	} else {
		data, err := rr.Func(r)
		if err.Message != "" {
			w.WriteHeader(err.Code)
			res = BuildJsonResponse(err)
		} else {
			w.WriteHeader(rr.StatusCode)
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
	r.Use(JwtTokenMiddleware)
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