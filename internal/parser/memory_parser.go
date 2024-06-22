package parser

import (
	"log"
	"time"
	"trustwallet/internal/storage"
	"trustwallet/internal/subscriber"
	"trustwallet/pkg/ethereum"
	"trustwallet/pkg/types"
)

type MemoryParser struct {
	subscriber *subscriber.Subscriber
}

func NewMemoryParser(url string) *MemoryParser {
	memoryStorage := storage.NewMemoryStorage()
	client := ethereum.NewDefaultClient(url)
	return &MemoryParser{subscriber: subscriber.NewSubscriber(memoryStorage, client)}
}

func (p *MemoryParser) Start() {
	go p.subscriber.Start()
}

func (p *MemoryParser) GetCurrentBlock() uint64 {
	return p.subscriber.GetLastParsedBlock()
}

func (p *MemoryParser) Subscribe(address string) bool {
	return p.subscriber.Subscribe(address)
}

func (p *MemoryParser) GetTransactions(address string) []types.Transaction {
	return p.subscriber.GetTransactions(address)
}

func (p *MemoryParser) RetrieveTransactionsPeriodically(address string, intervalSeconds int) {
	// Sample usage on how notification service can retrieve new incoming / outgoing transactions
	log.Printf("Initiated a periodic poll on Address [%s] at interval [%d seconds]", address, intervalSeconds)
	for {
		txns := p.GetTransactions(address)
		log.Printf("Transactions count retrieved [%d] Transactions %v Address [%s]\n", len(txns), txns, address)
		time.Sleep(time.Second * time.Duration(intervalSeconds))
	}
}