package http

import "net/http"

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type defaultClient struct{}

func (H defaultClient) Do(r *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(r)
}

func DefaultClient() Client {
	return defaultClient{}
}

type mockClient struct {
	// DoFunc will be executed whenever Do function is executed
	// so we'll be able to create a custom response
	DoFunc func(*http.Request) (*http.Response, error)
}

func (H mockClient) Do(r *http.Request) (*http.Response, error) {
	return H.DoFunc(r)
}

func MockClient() Client {
	return mockClient{}
}
