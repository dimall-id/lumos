package http

import (
	"net/http"
	"testing"
)

func TestRoute_IsValid(t *testing.T) {
	routes := []struct {
		route Route
		expect bool
	}{
		{
			route: Route{
				Name: "Product List",
				HttpMethod: "GET",
				Url: "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			expect: true,
		},
		{
			route: Route{
				HttpMethod: "GET",
				Url: "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			expect: false,
		},
		{
			route: Route{
				Name: "Product List",
				HttpMethod: "GETS",
				Url: "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			expect: false,
		},
		{
			route: Route{
				Name: "Product List",
				HttpMethod: "GETS",
				Url: "/pro du cts",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			expect: false,
		},
		{
			route: Route{
				Name: "Product List",
				HttpMethod: "GETS",
				Url: "/products",
			},
			expect: false,
		},
		{
			route: Route{
				Url: "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			expect: false,
		},
		{
			route: Route{
				Name: "Product List",
				HttpMethod: "POST",
				Url: "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			expect: true,
		},
	}

	for i, test := range routes {
		isValid := test.route.IsValid()
		if isValid != test.expect {
			isFunc := "PARSED"
			if test.route.Func == nil {
				isFunc = "NOT PARSED"
			}
			t.Errorf("[%d] Route testing is invalid HTTP METHOD : %s, URL : %s, NAME : %s, FUNC : %s", i, test.route.HttpMethod, test.route.Url, test.route.Name, isFunc)
		}
	}
}

func TestRoute_Equal(t *testing.T) {
	r1 := Route{
		HttpMethod: "GET",
		Url: "/products",
	}
	r2 := Route{
		HttpMethod: "GET",
		Url: "/products",
	}
	r3 := Route{
		HttpMethod: "POST",
		Url: "/products",
	}
	if !r1.Equal(r2) {
		t.Error("Fail to test, this should be equal")
	}
	if r1.Equal(r3) {
		t.Error("Fail to test, this should not be equal")
	}
}

func TestAddRoute(t *testing.T) {
	t.Run("Add Route", func(t *testing.T) {
		r1 := Route {
			Name: "Product List",
			HttpMethod: "GET",
			Url: "/products",
		}
		if err := AddRoute(r1); err == nil {
			t.Error("Fail to Test, Invalid Route still got added without error")
		}
		r2 := Route {
			Name: "Product List",
			HttpMethod: "GET",
			Url: "/products",
			Func: func(r *http.Request) (interface{}, HttpError) {
				return nil, HttpError{}
			},
		}
		if err := AddRoute(r2); err != nil {
			t.Error("Fail to test, an invalid route can't be added to routes")
		}
		if err := AddRoute(r2); err == nil {
			t.Error("Fail to test, an existing route can be added to routes")
		}
	})
}

func TestAddAllRoute(t *testing.T) {
	t.Run("Add All Route", func(t *testing.T) {

		r1s := []Route{
			{
				Name: "Product List",
				HttpMethod: "GET",
				Url: "/products",
			},
			{
				Name: "Product List",
				HttpMethod: "GET",
				Url: "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
		}
		if err := AddAllRoute(r1s); err == nil {
			t.Error("Fail to test, list contains of an invalid route. list can't be added to routes")
		}

		r2s := []Route{
			{
				Name: "Product List",
				HttpMethod: "GET",
				Url: "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			{
				Name: "Product List",
				HttpMethod: "GET",
				Url: "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
		}
		if err := AddAllRoute(r2s); err == nil {
			t.Error("Fail to test, list contain same route. Can't add same routes")
		}

		AddRoute(Route{
			Name: "Product List",
			HttpMethod: "GET",
			Url: "/products",
			Func: func(r *http.Request) (interface{}, HttpError) {
				return nil, HttpError{}
			},
		})

		r3s := []Route{
			{
				Name:       "Product List",
				HttpMethod: "GET",
				Url:        "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
		}
		if err := AddAllRoute(r3s); err == nil {
			t.Error("Fail to test, list contain of existing route.")
		}

		r4s := []Route{
			{
				Name:       "Product Add",
				HttpMethod: "POST",
				Url:        "/products",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			{
				Name:       "Product Delete",
				HttpMethod: "DELETE",
				Url:        "/products/{id}",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
			{
				Name:       "Product Update",
				HttpMethod: "PUT",
				Url:        "/products/{id}",
				Func: func(r *http.Request) (interface{}, HttpError) {
					return nil, HttpError{}
				},
			},
		}
		if err := AddAllRoute(r4s); err != nil {
			t.Errorf("Fail to test, this should success to add route. %s", err)
		}
	})
}

func TestGetRoute(t *testing.T) {
	r4s := []Route{
		{
			Name:       "Product Add",
			HttpMethod: "POST",
			Url:        "/products",
			Func: func(r *http.Request) (interface{}, HttpError) {
				return nil, HttpError{}
			},
		},
		{
			Name:       "Product Delete",
			HttpMethod: "DELETE",
			Url:        "/products/{id}",
			Func: func(r *http.Request) (interface{}, HttpError) {
				return nil, HttpError{}
			},
		},
		{
			Name:       "Product Update",
			HttpMethod: "PUT",
			Url:        "/products/{id}",
			Func: func(r *http.Request) (interface{}, HttpError) {
				return nil, HttpError{}
			},
		},
	}
	AddAllRoute(r4s)

	if r := GetRoute("POST", "/products"); r.HttpMethod != "POST" {
		t.Errorf("Fail to test, Fetching the wrong route")
	}

	r := GetRoute("GET", "/produts")
	if r.HttpMethod != "" {
		t.Errorf("Fail to Test, Successfully Fetch unexist route")
	}
}

func TestGetRouteAt(t *testing.T) {
	l := len(routes)
	r1 := GetRouteAt(l)
	if r1.HttpMethod != "" {
		t.Errorf("Fail to test, it should be returning empty route")
	}

	r2 := GetRouteAt(0)
	if r2.HttpMethod == "" {
		t.Errorf("Fail to test, it should not returning empty route")
	}
}

func TestGetAllRoute(t *testing.T) {
	rs := GetAllRoute()
	if len(rs) <= 0 {
		t.Errorf("Fail to test, it should not returning empty list")
	}
}