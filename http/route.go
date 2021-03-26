package http

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int
	Body interface{}
}

type Route struct {
	Name string
	HttpMethod string
	Url	string
	Roles []string
	Func func (r *http.Request) Response
}

func (r *Route) toFieldMaps() logrus.Fields {
	return map[string]interface{} {
		"name" : r.Name,
		"http_method" : r.HttpMethod,
		"url" : r.Url,
	}
}

func (r *Route) IsValid() error {
	if &r.Name == nil {
		return errors.New("name of route is not provided")
	}
	if &r.HttpMethod == nil {
		return errors.New("http method of route is not provided")
	}
	hm := "GETPOSTPUTPATCHDELETEOPTIONS"
	if !strings.Contains(hm, r.HttpMethod) {
		return errors.New(fmt.Sprintf("invalid http method provided, '%s'", r.HttpMethod))
	}
	if &r.Url == nil {
		return errors.New("url of route is not provided")
	}
	if &r.Func == nil {
		return errors.New("func of route is not provided")
	}
	if &r.Roles == nil {
		return errors.New("role of route is not provided")
	}
	return nil
}

func (r *Route) Equal (r2 Route) bool {
	if r.HttpMethod == r2.HttpMethod && r.Url == r2.Url {
		return true
	}
	return false
}

var routes []Route

func isExist (route Route) (bool, int) {
	for i, r := range routes {
		if r.Equal(route) {
			return true, i
		}
	}
	return false, -1
}

func AddRoute (route Route) error {
	log.Info("checking existing route")
	if oke,_ := isExist(route); oke {
		log.Warn("existing routes has been found")
		return &ExistingRouteError{route: route}
	}
	log.Info("checking route is valid")
	if err := route.IsValid(); err != nil {
		log.Errorf("route is not valid due to %s", err.Error())
		return &InvalidRouteError{route: route}
	} else {
		log.Info("appending new route")
		routes = append(routes, route)
		return nil
	}
}

func AddAllRoute (rs []Route) error {
	log.Info("checking if double route exist")
	for i, r1 := range rs {
		for j, r2 := range rs {
			if i != j && r1.Equal(r2) {
				return &DoubleRouteError{r1: i, r2: j}
			}
		}
	}
	for _,route := range rs {
		log.Info("checking existing route")
		if oke,_ := isExist(route); oke {
			log.Warn("existing routes has been found")
			return &ExistingRouteError{route: route}
		}
		log.Info("checking route is valid")
		if err := route.IsValid(); err != nil {
			log.Errorf("route is not valid due to %s", err.Error())
			return &InvalidRouteError{route: route}
		}
	}
	log.Info("appending new routes")
	routes = append(routes, rs...)
	return nil
}

func GetRoute (method string, url string) Route {
	route := Route{
		HttpMethod: method,
		Url: url,
		Func: nil,
	}
	if oke, ind := isExist(route); oke {
		return routes[ind]
	}
	return Route{}
}

func GetRouteAt (i int) Route {
	if i >= len(routes) {
		return Route{}
	}
	return routes[i]
}

func GetAllRoute () []Route {
	return routes
}