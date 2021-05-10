package main

import "os"

func main() {
	c := Client{}
	c.Run(os.Getenv("TRANSACTION_POOL_BASE_URL"))
}
