package http

import (
	"net/http"
	"strings"
	"testing"
)

func TestInternalServerError(t *testing.T) {
	err := InternalServerError("")
	if err.StatusCode != http.StatusInternalServerError {
		t.Error(err.StatusCode)
	}
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("")
	if err.StatusCode != http.StatusBadRequest {
		t.Error(err.StatusCode)
	}
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("")
	if err.StatusCode != http.StatusUnauthorized {
		t.Error(err.StatusCode)
	}
}

func TestPaymentRequired(t *testing.T) {
	err := PaymentRequired()
	if err.StatusCode != http.StatusPaymentRequired {
		t.Error(err.StatusCode)
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden("")
	if err.StatusCode != http.StatusForbidden {
		t.Error(err.StatusCode)
	}
}

func TestMethodNotAllow(t *testing.T) {
	err := MethodNotAllow("")
	if err.StatusCode != http.StatusMethodNotAllowed {
		t.Error(err.StatusCode)
	}
}

func TestUnprocessableEntity(t *testing.T) {
	err := UnprocessableEntity(map[string][]string{
		"username":[]string{
			"this is required",
			"that is required",
		},
	})
	if err.StatusCode != http.StatusUnprocessableEntity {
		t.Error(err.StatusCode)
	}
}

func TestInvalidRouteError_Error(t *testing.T) {
	r := Route{
		Name: "Name Product",
		HttpMethod: "GET",
		Url: "/products",
		Func: func(r *http.Request) Response {
			return Response{}
		},
	}
	err := InvalidRouteError{
		route: r,
	}
	if !strings.Contains(err.Error(), "Route given is invalid") {
		t.Errorf("Invalid Route Error")
	}
}

func TestExistingRouteError_Error(t *testing.T) {
	r := Route{
		Name: "Name Product",
		HttpMethod: "GET",
		Url: "/products",
		Func: func(r *http.Request) Response {
			return Response{}
		},
	}
	err := ExistingRouteError{
		route: r,
	}
	if !strings.Contains(err.Error(), "Existing Route with Http Method") {
		t.Errorf("Invalid Existing Route Error")
	}
}

func TestDoubleRouteError_Error(t *testing.T) {
	err := DoubleRouteError{
		r1: 1,
		r2: 2,
	}
	if !strings.Contains(err.Error(), "Double Route found in the list") {
		t.Errorf("Invalid Double Route")
	}
}