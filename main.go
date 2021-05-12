package main

import (
	"net/http"
	"os"
)

func main() {
	baseUrl := os.Getenv("TRANSACTION_POOL_BASE_URL")
	c := TransactionSpawner{&http.Client{}, baseUrl}
	c.SpawnTransaction()
}
