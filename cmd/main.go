package main

import (
	"time"
	"trustwallet/internal/parser"
)

func main() {
	p := parser.NewMemoryParser("https://cloudflare-eth.com")
	subscriptionAddress := "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5" // beaverbuild block builder address
	p.Start()
	p.Subscribe(subscriptionAddress)
	go p.RetrieveTransactionsPeriodically(subscriptionAddress, 60) // sample call on how notification service can fetch new transactions

	block() // make a blocking call to log out and demonstrate functionality of this library
}

func block() {
	for {
		time.Sleep(time.Second)
	}
}
