package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dimall-id/jwt-go"
	"github.com/dimall-id/lumos/v2/misc"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var _publicKey []byte

func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := MethodNotAllow("method is not allowed")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(err.StatusCode)
		w.Write(BuildJsonResponse(err.Body))
	})
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := NotFound("url is not found")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(err.StatusCode)
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
	log.Infoln("converting array of role to may of roles")
	roles := make(map[string]string)
	for _, role := range rr.Roles {roles[role] = role}
	log.Infof("checking the len of roles in routes. len = %d", len(roles))
	if len(roles) <= 0 {return Response{}}
	log.Infof("checking if roles has ANONYMOUS role")
	if _, oke := roles["ANONYMOUS"]; oke {
		log.Infof("routes contains ANONYMOUS role")
		return Response{}
	}

	if authentication == "" {
		log.WithField("User-Id", "")
		log.Infof("authorization key if not provided in the header")
		return Unauthorized("authorization key is not provided in the header")
	} else {
		log.Infof("parsing and validating the jwt token signature")
		tokens := misc.BuildToMap(`Bearer (?P<token>[\W\w]+)`, authentication)
		claims, err := jwt.ParseWithClaims(tokens["token"], jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {return _publicKey, nil})
		if err != nil {
			log.Infof("fail to validate the jwt token signature")
			log.Error(err)
			vErr := err.(*jwt.ValidationError)
			return Forbidden(vErr.Error())
		}
		log.Infof("parsing token claim to AcessToken")
		accessToken := AccessToken{}
		accessToken.FillAccessToken(claims.Claims.(jwt.MapClaims))
		log.WithField("User-Id", accessToken.UserId)
		log.Infof("checking issued at and expired at")
		err = accessToken.Valid()
		if err != nil {
			log.Errorf("invalid access token due to %s", err.Error())
			return Forbidden(fmt.Sprintf("invalid access token due to %s", err.Error()))
		}
		log.Infof("checking role of user vs route role")
		if !CheckRole(accessToken.Roles, roles) {
			return Unauthorized("user don't have role to access the resources")
		}
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
	reqId := uuid.New().String()
	r.WithContext(context.WithValue(r.Context(), "X-Req-Id", reqId))
	log.WithField("X-Req-Id", reqId)
	var res []byte
	log.Infof("start handling request for url %s", r.RequestURI)
	log.Infoln("checking the authorization")
	err := CheckAuthorization(r.Header.Get("Authorization"), rr)
	if err.StatusCode != 0 {
		log.Infoln("fail to check authorization")
		w.WriteHeader(err.StatusCode)
		res = BuildJsonResponse(err.Body)
	} else {
		resp := rr.Func(r)
		log.Infof("process request return with status code %d\n", resp.StatusCode)
		if resp.StatusCode == 0 {w.WriteHeader(http.StatusOK)} else {w.WriteHeader(resp.StatusCode)}
		res = BuildJsonResponse(resp.Body)
	}
	log.WithField("Response Size", len(res)).Infof("done handling request for url \"%s\"", r.RequestURI)
	w.Write(res)
}

func GenerateMuxRouter (routes []Route, middleware []mux.MiddlewareFunc) *mux.Router {
	log.Infoln("initializing mux router")
	r := mux.NewRouter()
	log.Infoln("registering method not found handler")
	r.MethodNotAllowedHandler = methodNotAllowedHandler()
	log.Infoln("registering not found handler")
	r.NotFoundHandler = notFoundHandler()

	log.Infof("register routes, total %d routes", len(routes))
	for i, _ := range routes {
		log.WithFields(routes[i].toFieldMaps()).Infof("registering routes %s", routes[i].Name)
		rr := GetRouteAt(i)
		r.HandleFunc(rr.Url, func(w http.ResponseWriter, r *http.Request) {
			HandleRequest(w, r, rr)
		}).Methods(rr.HttpMethod).Name(rr.Name)
	}

	log.Info("registering content type middleware")
	r.Use(ContentTypeMiddleware)
	log.Info("registering jwt middleware")
	r.Use(JwtTokenMiddleware)
	log.Infof("registering middlewares, total %d middlewares", len(middleware))
	for _, mwr := range middleware {
		mw := mwr
		r.Use(mw)
	}

	return r
}

func setPublicKey (publicKeyUrl string) error {
	log.Infoln("checking if public key url exists")
	if publicKeyUrl == "" {
		return errors.New("public key url not provided")
	}

	log.Infoln("fetching the public key content")
	resp, err := http.Get(publicKeyUrl)
	if err != nil {
		log.Errorln(err)
		return errors.New("fail to consume public key url")
	}

	log.Infoln("reading public key response content")
	byteKey, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return errors.New("fail to read http response")
	}

	_publicKey = byteKey
	return nil
}

func StartHttpServer(port string, publicKey string) error {
	err := setPublicKey(publicKey)
	if err != nil {return err}
	r := GenerateMuxRouter(routes, middlewares)
	return http.ListenAndServe(port, r)
}