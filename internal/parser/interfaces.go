package parser

import "trustwallet/pkg/types"

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// Start all relevant services
	Start() error

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []types.Transaction
}
