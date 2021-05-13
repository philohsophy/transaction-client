package mocks

import (
	"net/http"
)

// Ref: https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/
type ClientMock struct{}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

// Ref: http://hassansin.github.io/Unit-Testing-http-client-in-Go
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
