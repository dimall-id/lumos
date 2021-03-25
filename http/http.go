package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dimall-id/jwt-go"
	"github.com/dimall-id/lumos/misc"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var _publicKey string

func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := MethodNotAllow("method is not allowed")
		w.WriteHeader(err.StatusCode)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Write(BuildJsonResponse(err.Body))
	})
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := NotFound("url is not found")
		w.WriteHeader(err.StatusCode)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Write(BuildJsonResponse(err.Body))
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
func CheckAuthorization(authentication string, rr Route) Response {
	roles := make(map[string]string)
	for _, role := range rr.Roles {roles[role] = role}
	if len(roles) <= 0 {return Response{}}
	if _, oke := roles["ANONYMOUS"]; oke {return Response{}}

	if authentication == "" {
		return Unauthorized("authorization key is not provided in the header")
	} else {
		var claims jwt.MapClaims
		tokens := misc.BuildToMap(`Bearer (?P<token>[\W\w]+)`, authentication)
		_, err := jwt.ParseWithClaims(tokens["token"], claims, func(token *jwt.Token) (interface{}, error) {return _publicKey, nil})
		if err != nil {
			vErr := err.(jwt.ValidationError)
			return Forbidden(vErr.Error())
		}
		accessToken := AccessToken{}
		accessToken.FillAccessToken(claims)
		err = accessToken.Valid()
		if err != nil {return Forbidden(err.Error())}
		if !CheckRole(accessToken.Roles, roles) {return Unauthorized("user don't have role to access the resources")}
	}
	return Response{}
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
	if &err != nil {
		w.WriteHeader(err.StatusCode)
		res = BuildJsonResponse(err)
	} else {
		resp := rr.Func(r)
		if resp.StatusCode != 0 {w.WriteHeader(http.StatusOK)} else {w.WriteHeader(resp.StatusCode)}
		res = BuildJsonResponse(resp.Body)
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

func setPublicKey (publicKeyUrl string) error {
	if publicKeyUrl == "" {return errors.New("public key url not provided")}

	resp, err := http.Get(publicKeyUrl)
	if err != nil {return errors.New("fail to consume public key url")}

	byte, err := ioutil.ReadAll(resp.Body)
	if err != nil {return errors.New("fail to read http response")}

	_publicKey = string(byte)
	return nil
}

func StartHttpServer(port string, publicKey string) error {
	err := setPublicKey(publicKey)
	if err != nil {return err}
	r := GenerateMuxRouter(routes, middlewares)
	return http.ListenAndServe(port, r)
}