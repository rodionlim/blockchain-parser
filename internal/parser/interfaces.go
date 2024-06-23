package parser

import (
	"context"
	"parser/pkg/types"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// Start parser and all relevant services
	Start(ctx context.Context, cancel context.CancelFunc) error

	// Stop parser and all relevant services
	Stop() error

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []types.Transaction
}
