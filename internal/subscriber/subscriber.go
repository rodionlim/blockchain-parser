package subscriber

import (
	"log"
	"strings"
	"sync"
	"time"
	"trustwallet/internal/storage"
	"trustwallet/pkg/ethereum"
	"trustwallet/pkg/types"
)

type Subscriber struct {
	storage                  storage.Storage
	client                   ethereum.Client
	lastParsedBlock          uint64
	pollingDurationInSeconds int
	mu                       sync.RWMutex
}

func NewSubscriber(storage storage.Storage, client ethereum.Client) *Subscriber {
	var bn uint64
	var err error
	if bn, err = client.GetBlockNumber(); err != nil {
		panic("Unable to connect to blockchain client")
	}
	log.Printf("New Subscriber created with last parsed block [%d]\n", bn)
	return &Subscriber{storage: storage, client: client, lastParsedBlock: bn, pollingDurationInSeconds: 15}
}

func (s *Subscriber) Start() {
	for {
		s.mu.Lock()
		results, err := s.client.GetTransactionsBySubscriptionsAndBlock(s.storage.GetAllSubscriptions(), s.lastParsedBlock)
		if err != nil {
			panic("Something went wrong while trying to retrieve transactions from a block")
		}
		s.storage.WriteTransactions(results)
		s.lastParsedBlock++
		s.mu.Unlock()
		bn, err := s.client.GetBlockNumber()
		if err != nil {
			panic("Unable to fetch latest block number")
		}
		if bn > s.lastParsedBlock {
			// skip the sleep if we are behind by too much
			log.Println("Latest block is ahead by > 1 polling interval, polling for the next block immediately")
			continue
		}
		time.Sleep(time.Duration(time.Second) * time.Duration(s.pollingDurationInSeconds))
	}
}

func (s *Subscriber) Subscribe(address string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	address = strings.ToLower(address)
	if !s.storage.AddressExists(address) {
		log.Printf("Made Subscription to Address [%s]\n", address)
		s.storage.SetFirstParsedBlock(address, s.lastParsedBlock+1)
		return true
	} else {
		// already subscribed
		return false
	}
}

func (s *Subscriber) GetSubscribedAddresses() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.storage.GetAllAddresses()
}

func (s *Subscriber) GetLastParsedBlock() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lastParsedBlock
}

func (s *Subscriber) GetTransactions(address string) []types.Transaction {
	data, err := s.storage.ReadTransactions(address)
	if err != nil {
		log.Panicf("Something went wrong while trying to read transactions for address [%s]", err)
	}
	return data
}
