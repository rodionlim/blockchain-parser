package parser

import (
	"context"
	"trustwallet/pkg/types"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// Start parser and all relevant services
	Start(ctx context.Context) error

	// Stop parser and all relevant services
	Stop(cancel context.CancelFunc) error

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []types.Transaction
}
