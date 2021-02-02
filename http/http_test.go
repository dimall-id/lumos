package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Todo struct {
	Val string `json:"val"`
	IsTick bool `json:"is_tick"`
}

func TestHandleRequest(t *testing.T) {
	t.Run("TESTING HANDLE REQUEST", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			HandleRequest(writer, request, func(r2 *http.Request) (interface{}, HttpError) {
				datas := []Todo {
					{
						Val: "Mencuci Pakaian",
						IsTick: true,
					},
				}
				return datas, HttpError{}
			})
		})

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Fail to test, the response should return 200 status")
		}

		ct := rr.Header().Get("Content-Type")
		if !strings.Contains(ct, "application/json") {
			t.Errorf("Fail to test, the response should have content type of application/json")
		}
	})
	t.Run("TESTING HANDLE REQUEST 2", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			HandleRequest(writer, request, func(r2 *http.Request) (interface{}, HttpError) {
				datas := []Todo {
					{
						Val: "Mencuci Pakaian",
						IsTick: true,
					},
				}
				return datas, NotFound()
			})
		})

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Fail to test, the response should return 404 status")
		}
	})
}