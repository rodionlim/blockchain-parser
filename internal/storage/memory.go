package storage

// in-memory storage implementation of the storage interface

import (
	"fmt"
	"log"
	"parser/pkg/types"
	"strings"
	"sync"
)

type MemoryStorage struct {
	trackerData map[string]*types.Tracker      // address to subscription tracking data mapping
	txnData     map[string][]types.Transaction // address to transactions mapping
	mu          sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		trackerData: make(map[string]*types.Tracker),
		txnData:     make(map[string][]types.Transaction),
		mu:          sync.RWMutex{},
	}
}

func (m *MemoryStorage) AddressExists(address string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.trackerData[address]
	return exists
}

func (m *MemoryStorage) SetFirstParsedBlock(address string, blockNumber uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	address = strings.ToLower(address)
	data, ok := m.trackerData[address]
	if ok {
		data.FirstParsedBlock = blockNumber
	} else {
		m.trackerData[address] = &types.Tracker{FirstParsedBlock: blockNumber, Offset: 0}
	}
}

func (m *MemoryStorage) GetAllAddresses() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	addresses := make([]string, 0, len(m.trackerData))
	for address := range m.trackerData {
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func (m *MemoryStorage) GetAllSubscriptions() map[string]*types.Tracker {
	return m.trackerData
}

func (m *MemoryStorage) WriteTransactions(txnsMap map[string][]types.Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for address := range txnsMap {
		address = strings.ToLower(address)
		txns := txnsMap[address]
		_, ok := m.txnData[address]
		if ok {
			m.txnData[address] = append(m.txnData[address], txns...)
		} else {
			m.txnData[address] = txns
		}
	}
}

func (m *MemoryStorage) ReadTransactions(address string) ([]types.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	address = strings.ToLower(address)
	metadata, ok := m.trackerData[address]
	if !ok {
		return nil, fmt.Errorf("invalid address [%s] specified. Address is currently not being subscribed to", address)
	}
	data, ok := m.txnData[address]
	if !ok {
		return nil, nil
	}
	log.Printf("Reading Transactions for Address [%s] Current Offset [%d] Total Txn Count [%d]\n", address, metadata.Offset, len(data))
	newOffset := len(data)
	data = data[metadata.Offset:]
	m.trackerData[address].Offset = newOffset // update the offset pointer so that we don't duplicate any notifications back to the client
	return data, nil
}
