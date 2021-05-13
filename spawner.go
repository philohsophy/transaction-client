package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/philohsophy/dummy-blockchain-models"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TransactionSpawner struct {
	Client  HTTPClient
	baseUrl string
}

func (ts *TransactionSpawner) SpawnTransaction() error {
	transaction := createTransaction()
	return sendTransaction(ts, transaction)
}

var createTransaction = func() models.Transaction {
	var t models.Transaction
	t.RecipientAddress = models.Address{Name: "Foo", Street: "FooStreet", HouseNumber: "1", Town: "FooTown"}
	t.SenderAddress = models.Address{Name: "Bar", Street: "BarStreet", HouseNumber: "1", Town: "BarTown"}
	t.Value = 100.21

	return t
}

var jsonMarshal = json.Marshal
var httpNewRequest = http.NewRequest

var sendTransaction = func(ts *TransactionSpawner, transaction models.Transaction) error {
	data, err := jsonMarshal(transaction)
	if err != nil {
		return fmt.Errorf("Error transforming transaction to JSON: %s", err.Error())
	}

	endpoint := ts.baseUrl + "/transactions"
	req, err := httpNewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := ts.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Error reading response: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		reason := res.Body
		return fmt.Errorf("Failed to insert Transaction into Transaction-Pool: %s", reason)
	}

	return nil
}
