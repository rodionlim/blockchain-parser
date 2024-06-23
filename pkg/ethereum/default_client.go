package ethereum

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"parser/pkg/types"
)

type DefaultClient struct {
	URL string
}

func NewDefaultClient(url string) *DefaultClient {
	return &DefaultClient{URL: url}
}

func (c *DefaultClient) Call(method string, params interface{}) (json.RawMessage, error) {
	reqBody, err := json.Marshal(types.RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp types.RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, err
	}

	if rpcResp.Error != nil {
		return nil, errors.New(rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

func (c *DefaultClient) GetBlockNumber() (uint64, error) {
	result, err := c.Call("eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, err
	}

	var blockNumberHex string
	if err := json.Unmarshal(result, &blockNumberHex); err != nil {
		return 0, err
	}

	var blockNumber uint64
	fmt.Sscanf(blockNumberHex, "0x%x", &blockNumber)
	return blockNumber, nil
}

func (c *DefaultClient) GetBlockByNumber(blockNumber uint64) (map[string]interface{}, error) {
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
	result, err := c.Call("eth_getBlockByNumber", []interface{}{blockNumberHex, true})
	if err != nil {
		return nil, err
	}

	var block map[string]interface{}
	if err := json.Unmarshal(result, &block); err != nil {
		return nil, err
	}

	return block, nil
}

func (c *DefaultClient) GetTransactionsBySubscriptionsAndBlock(subscriptions map[string]*types.Tracker, blockNumber uint64) (map[string][]types.Transaction, error) {
	// returns a map of address to list of inbound/outbound transactions
	results := make(map[string][]types.Transaction)
	if len(subscriptions) == 0 {
		log.Printf("Current Block Number [%d] No subscriptions found\n", blockNumber)
		return results, nil
	}
	block, err := c.GetBlockByNumber(blockNumber)
	if err != nil {
		return nil, err
	}

	count := 0
	found := 0
	for _, tx := range block["transactions"].([]interface{}) {
		count++
		txMap := tx.(map[string]interface{})
		from, ok := txMap["from"].(string)
		if !ok {
			from = ""
		}
		to, ok := txMap["to"].(string)
		if !ok {
			from = ""
		}
		_, fromExists := subscriptions[from]
		_, toExists := subscriptions[to]
		if fromExists {
			txns, ok := results[from]
			if ok {
				results[from] = append(txns, txMap)
			} else {
				results[from] = []types.Transaction{txMap}
			}
		}
		if toExists {
			txns, ok := results[to]
			if ok {
				results[to] = append(txns, txMap)
			} else {
				results[to] = []types.Transaction{txMap}
			}
		}
		if fromExists || toExists {
			found++
		}
	}
	log.Printf("Queried Block Number [%d] Total transactions [%d] Relevant transactions from subscriptions count [%d]\n", blockNumber, count, found)
	return results, nil
}
