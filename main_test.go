package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var baseUrl string = "http://localhost:8011"

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

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

// equals fails the test if exp is not equal to act.
// Ref: https://github.com/benbjohnson/testing
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
// Ref: https://github.com/benbjohnson/testing
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func TestSpawnTransaction(t *testing.T) {

	t.Run("If server accepts the transaction", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			equals(t, req.URL.String(), baseUrl+"/transactions")
			return &http.Response{
				StatusCode: 201,
				Header:     make(http.Header),
			}
		})

		ts := TransactionSpawner{client, baseUrl}
		err := ts.SpawnTransaction()
		ok(t, err)
	})

	t.Run("If server does not accept the transaction (4xx)", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			equals(t, req.URL.String(), baseUrl+"/transactions")
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewBufferString(`Error Message`)),
				Header:     make(http.Header),
			}
		})

		ts := TransactionSpawner{client, baseUrl}
		err := ts.SpawnTransaction()
		if err == nil {
			t.Error("Expected SpawnTransaction to return an error")
		}
	})

	t.Run("If server a server side error occures", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			equals(t, req.URL.String(), baseUrl+"/transactions")
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(`Server Error Message`)),
				Header:     make(http.Header),
			}
		})

		ts := TransactionSpawner{client, baseUrl}
		err := ts.SpawnTransaction()
		if err == nil {
			t.Error("Expected SpawnTransaction to return an error")
		}
	})
}
