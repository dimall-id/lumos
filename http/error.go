package http

import "net/http"

type HttpError struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Errors map[string][]string `json:"errors,omitempty"`
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