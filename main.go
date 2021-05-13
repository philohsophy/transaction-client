package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	baseUrl := os.Getenv("TRANSACTION_POOL_BASE_URL")
	c := TransactionSpawner{&http.Client{}, baseUrl}
	err := c.SpawnTransaction()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
