package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	models "github.com/philohsophy/dummy-blockchain-models"
)

type Client struct{}

func (c *Client) Run(baseUrl string) {
	t := createTransaction()
	sendTransaction(t, baseUrl)
}

func createTransaction() models.Transaction {
	var t models.Transaction
	t.RecipientAddress = models.Address{Name: "Foo", Street: "FooStreet", HouseNumber: "1", Town: "FooTown"}
	t.SenderAddress = models.Address{Name: "Bar", Street: "BarStreet", HouseNumber: "1", Town: "BarTown"}
	t.Value = 100.21

	return t
}

func sendTransaction(transaction models.Transaction, baseUrl string) {
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Fatal("Error transforming transaction to JSON", err)
	}

	endpoint := baseUrl + "/transactions"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	fmt.Printf("%s\n", body)
}
