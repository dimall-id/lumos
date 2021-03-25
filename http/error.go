package http

import (
	"fmt"
	"net/http"
)

func InternalServerError(cause string) Response {
	return Response{
		StatusCode: http.StatusInternalServerError,
		Body: map[string]interface{} {
			"message" : "internal server error",
			"cause" : cause,
		},
	}
}

func BadRequest(cause string) Response {
	return Response {
		StatusCode: http.StatusBadRequest,
		Body: map[string]interface{} {
			"message" : "bad request",
			"cause" : cause,
		},
	}
}

func Unauthorized(cause string) Response {
	return Response {
		StatusCode: http.StatusUnauthorized,
		Body: map[string]interface{} {
			"message" : "unauthorized",
			"cause" : cause,
		},
	}
}

func PaymentRequired() Response {
	return Response {
		StatusCode: http.StatusPaymentRequired,
		Body: map[string]interface{} {
			"message" : "payment required",
		},
	}
}

func Forbidden(cause string) Response {
	return Response {
		StatusCode: http.StatusForbidden,
		Body: map[string]interface{} {
			"message" : "forbidden",
			"cause" : cause,
		},
	}
}

func NotFound(cause string) Response {
	return Response {
		StatusCode: http.StatusNotFound,
		Body: map[string]interface{} {
			"message" : "not found",
			"cause" : cause,
		},
	}
}

func MethodNotAllow(cause string) Response {
	return Response {
		StatusCode: http.StatusMethodNotAllowed,
		Body: map[string]interface{} {
			"message" : "method not allowed",
			"cause" : cause,
		},
	}
}

func UnprocessableEntity(errors map[string][]string) Response {
	return Response{
		StatusCode: http.StatusUnprocessableEntity,
		Body: map[string]interface{} {
			"message" : "unprocessable entity",
			"errors" : errors,
		},
	}
}

func NotImplemented(cause string) Response {
	return Response{
		StatusCode: http.StatusNotImplemented,
		Body: map[string]interface{} {
			"message" : "unprocessable entity",
			"cause" : cause,
		},
	}
}

type InvalidRouteError struct {
	route Route
}

func (ir *InvalidRouteError) Error() string {
	return "Route given is invalid"
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