package ethereum

import (
	"parser/pkg/types"
)

type Client interface {
	GetBlockNumber() (uint64, error)                                                                                                            // retrieve latest block number
	GetBlockByNumber(blockNumber uint64) (map[string]interface{}, error)                                                                        // retrieve block data by block number
	GetTransactionsBySubscriptionsAndBlock(subscriptions map[string]*types.Tracker, blockNumber uint64) (map[string][]types.Transaction, error) // retrieve relevant transactions by block number
}
