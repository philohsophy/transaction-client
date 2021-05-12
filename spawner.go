package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	models "github.com/philohsophy/dummy-blockchain-models"
)

type TransactionSpawner struct {
	Client  *http.Client
	baseUrl string
}

func (ts *TransactionSpawner) SpawnTransaction() error {
	t := ts.createTransaction()
	return ts.sendTransaction(t)
}

func (ts *TransactionSpawner) createTransaction() models.Transaction {
	var t models.Transaction
	t.RecipientAddress = models.Address{Name: "Foo", Street: "FooStreet", HouseNumber: "1", Town: "FooTown"}
	t.SenderAddress = models.Address{Name: "Bar", Street: "BarStreet", HouseNumber: "1", Town: "BarTown"}
	t.Value = 100.21

	return t
}

func (ts *TransactionSpawner) sendTransaction(transaction models.Transaction) error {
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Fatal("Error transforming transaction to JSON", err)
	}

	endpoint := ts.baseUrl + "/transactions"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := ts.Client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		reason := res.Body
		return fmt.Errorf("Failed to create Transaction: %s", reason)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	fmt.Printf("%s\n", body)

	return nil
}
