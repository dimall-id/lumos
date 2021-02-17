package http

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	Code int `json:"code" msgpack:"code"`
	Message string `json:"message" msgpack:"message"`
	Errors map[string][]string `json:"errors,omitempty" msgpack:"errors,as_array,omitempty"`
}

func InternalServerError() HttpError {
	return HttpError{
		Code: http.StatusInternalServerError,
		Message: "Internal Server Error",
		Errors: nil,
	}
}

func BadRequest() HttpError {
	return HttpError{
		Code: http.StatusBadRequest,
		Message: "Bad Request",
		Errors: nil,
	}
}

func Unauthorized() HttpError {
	return HttpError{
		Code: http.StatusUnauthorized,
		Message: "Unauthorized",
		Errors: nil,
	}
}

func PaymentRequired() HttpError {
	return HttpError{
		Code: http.StatusPaymentRequired,
		Message: "Payment Required",
		Errors: nil,
	}
}

func Forbidden() HttpError {
	return HttpError{
		Code: http.StatusForbidden,
		Message: "Forbidden",
		Errors: nil,
	}
}

func NotFound() HttpError {
	return HttpError{
		Code: http.StatusNotFound,
		Message: "Not Found",
		Errors: nil,
	}
}

func MethodNotAllow() HttpError {
	return HttpError{
		Code: http.StatusMethodNotAllowed,
		Message: "Method Not Allowed",
		Errors: nil,
	}
}

func UnprocessableEntity(errors map[string][]string) HttpError {
	return HttpError{
		Code: http.StatusUnprocessableEntity,
		Message: "Unprocessable Entity",
		Errors: errors,
	}
}

type InvalidRouteError struct {
	route Route
}

func (ir *InvalidRouteError) Error() string {
	var funcStatus string
	if ir.route.Func == nil {
		funcStatus = "NOT PARSED"
	} else {
		funcStatus = "PARSED"
	}
	return fmt.Sprintf("Route given is invalid. Http Method : %s, Url : %s, Func : %s", ir.route.HttpMethod, ir.route.Url, funcStatus)
}

type ExistingRouteError struct {
	route Route
}

func (er *ExistingRouteError) Error() string {
	return fmt.Sprintf("Existing Route with Http Method : %s and Url : %s existed", er.route.HttpMethod, er.route.Url)
}

type DoubleRouteError struct {
	r1 int
	r2 int
}

func (dr *DoubleRouteError) Error() string {
	return fmt.Sprintf("Double Route found in the list. Check at index %d and %d", dr.r1, dr.r2)
}