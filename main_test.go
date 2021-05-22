package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/google/uuid"
	models "github.com/philohsophy/blockchain-models"
	"github.com/philohsophy/transaction-spawner/mocks"
)

var baseUrl string = "http://localhost:8011"

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
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

func TestCreateTransaction(t *testing.T) {
	transaction := createTransaction()
	if transaction.IsValid() {
		t.Error("Expected created Transaction to be invalid without an id")
	}
	transaction.Id = uuid.New()
	if !transaction.IsValid() {
		t.Error("Expected created Transaction to be valid with an id")
	}
}

func TestSendTransaction(t *testing.T) {

	t.Run("If the server accepts the transaction", func(t *testing.T) {
		client := mocks.NewTestClient(func(req *http.Request) *http.Response {
			equals(t, req.URL.String(), baseUrl+"/transactions")
			return &http.Response{
				StatusCode: 201,
				Header:     make(http.Header),
			}
		})

		ts := TransactionSpawner{client, baseUrl}
		transaction := createTransaction()
		err := sendTransaction(&ts, transaction)
		ok(t, err)
	})

	t.Run("If the transaction cannot be converted to JSON", func(t *testing.T) {
		// mock jsonMarshal (preserve original function)
		jsonMarshalOrig := jsonMarshal
		jsonMarshal = func(v interface{}) ([]byte, error) {
			return nil, fmt.Errorf("jsonMarshalError")
		}

		ts := TransactionSpawner{nil, baseUrl}
		transaction := createTransaction()
		err := sendTransaction(&ts, transaction)
		if err == nil {
			t.Errorf("Expected sendTransaction to return an error. Got '%v'", err)
		}

		// revert jsonMarshal
		jsonMarshal = jsonMarshalOrig
	})

	t.Run("If the http request cannot be createed", func(t *testing.T) {
		// mock httpNewRequest (preserve original function)
		httpNewRequestOrig := httpNewRequest
		httpNewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
			return nil, fmt.Errorf("httpNewRequestError")
		}

		ts := TransactionSpawner{nil, baseUrl}
		transaction := createTransaction()
		err := sendTransaction(&ts, transaction)
		if err == nil {
			t.Errorf("Expected sendTransaction to return an error'. Got '%v'", err)
		}

		// revert httpNewRequest
		httpNewRequest = httpNewRequestOrig
	})

	t.Run("If performing the http request fails (client.Do())", func(t *testing.T) {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return nil, errors.New(
				"httpDoRequestError",
			)
		}

		ts := TransactionSpawner{&mocks.ClientMock{}, baseUrl}
		transaction := createTransaction()
		err := sendTransaction(&ts, transaction)
		if err == nil {
			t.Errorf("Expected sendTransaction to return an error'. Got '%v'", err)
		}
	})

	t.Run("If server does not return a 201", func(t *testing.T) {
		statusCodes := [2]int{400, 500}

		for _, statusCode := range statusCodes {
			client := mocks.NewTestClient(func(req *http.Request) *http.Response {
				equals(t, req.URL.String(), baseUrl+"/transactions")
				return &http.Response{
					StatusCode: statusCode,
					Body:       io.NopCloser(bytes.NewBufferString(`Error Message`)),
					Header:     make(http.Header),
				}
			})

			ts := TransactionSpawner{client, baseUrl}
			err := ts.SpawnTransaction()
			if err == nil {
				t.Errorf("Expected sendTransaction to return an error'. Got '%v'", err)
			}
		}
	})
}

func TestSpawnTransaction(t *testing.T) {
	t.Run("If everything runs fine", func(t *testing.T) {
		client := mocks.NewTestClient(func(req *http.Request) *http.Response {
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

	t.Run("If sendTransaction returns an error", func(t *testing.T) {
		sendTransactionError := "sendTransactionError"
		// mock sendTransaction (preserve original function)
		sendTransactionOrig := sendTransaction
		sendTransaction = func(ts *TransactionSpawner, transaction models.Transaction) error {
			return fmt.Errorf(sendTransactionError)
		}

		ts := TransactionSpawner{nil, baseUrl}
		err := ts.SpawnTransaction()
		if err == nil {
			t.Errorf("Expected SpawnTransaction to return an error '%s'. Got '%v'", sendTransactionError, err)
		}

		// revert sendTransaction
		sendTransaction = sendTransactionOrig
	})
}
