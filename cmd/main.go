package main

import (
	"context"
	"log"
	"parser/internal/parser"
	"time"
)

func main() {
	p := parser.NewMemoryParser("https://cloudflare-eth.com")
	subscriptionAddress := "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5" // beaverbuild block builder address
	ctx, cancel := context.WithCancel(context.Background())

	p.Start(ctx, cancel)                                           // start parser service
	p.Subscribe(subscriptionAddress)                               // subscribe to new transactions arising from beaverbuild block builder
	go p.RetrieveTransactionsPeriodically(subscriptionAddress, 60) // sample call on how notification service can fetch new transactions

	block(time.Minute * 3) // make a blocking call to log out and demonstrate functionality of this library
	p.Stop()               // perform all necessary clean up activities
}

func block(blockingPeriod time.Duration) {
	log.Printf("Running this application for Period [%v]", blockingPeriod.String())
	startTime := time.Now()
	for {
		duration := time.Since(startTime)
		if duration > blockingPeriod {
			return
		}
		time.Sleep(time.Second)
	}
}
