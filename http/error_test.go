package http

import (
	"net/http"
	"strings"
	"testing"
)

func TestInternalServerError(t *testing.T) {
	err := InternalServerError()
	if err.Code != http.StatusInternalServerError {
		t.Error(err.Code, err.Message)
	}
}

func TestBadRequest(t *testing.T) {
	err := BadRequest()
	if err.Code != http.StatusBadRequest {
		t.Error(err.Code, err.Message)
	}
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized()
	if err.Code != http.StatusUnauthorized {
		t.Error(err.Code, err.Message)
	}
}

func TestPaymentRequired(t *testing.T) {
	err := PaymentRequired()
	if err.Code != http.StatusPaymentRequired {
		t.Error(err.Code, err.Message)
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden()
	if err.Code != http.StatusForbidden {
		t.Error(err.Code, err.Message)
	}
}

func TestMethodNotAllow(t *testing.T) {
	err := MethodNotAllow()
	if err.Code != http.StatusMethodNotAllowed {
		t.Error(err.Code, err.Message)
	}
}

func TestUnprocessableEntity(t *testing.T) {
	err := UnprocessableEntity(map[string][]string{
		"username": []string{
			"this is required",
			"that is required",
		},
	})
	if err.Code != http.StatusUnprocessableEntity {
		t.Error(err.Code, err.Message)
	}
	if len(err.Errors) != 1 {
		t.Errorf("Errors detail is not set")
	}
	if len(err.Errors["username"]) != 2 {
		t.Errorf("Field Error Detail is not set")
	}
}

func TestInvalidRouteError_Error(t *testing.T) {
	//r := Route{
	//	Name:       "Name Product",
	//	HttpMethod: "GET",
	//	Url:        "/products",
	//	Func: func(r *http.Request) (interface{}, HttpError) {
	//		return nil, HttpError{}
	//	},
	//}
	err := InvalidRouteError{
		msg: "Route given is invalid",
	}
	if !strings.Contains(err.Error(), "Route given is invalid") {
		t.Errorf("Invalid Route Error")
	}
}

func TestExistingRouteError_Error(t *testing.T) {
	r := Route{
		Name:       "Name Product",
		HttpMethod: "GET",
		Url:        "/products",
		Func: func(r *http.Request) (interface{}, HttpError) {
			return nil, HttpError{}
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
