package storage

import "parser/pkg/types"

type Storage interface {
	AddressExists(address string) bool                            // check if an address is being watched
	SetFirstParsedBlock(address string, blockNumber uint64)       // keep a record of when an address is added as an observer
	GetAllAddresses() ([]string, error)                           // fetch a list of all subscriptions addresses
	GetAllSubscriptions() map[string]*types.Tracker               // fetch a list of all subscriptions and their corresponding metadata
	WriteTransactions(txnsMap map[string][]types.Transaction)     // persist transactions by addresses
	ReadTransactions(address string) ([]types.Transaction, error) // read transactions starting from offset and update new offset
}
