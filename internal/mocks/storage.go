package mocks

import (
	"parser/pkg/types"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetAllSubscriptions() map[string]*types.Tracker {
	args := m.Called()
	return args.Get(0).(map[string]*types.Tracker)
}

func (m *MockStorage) WriteTransactions(txs map[string][]types.Transaction) {
	m.Called(txs)
}

func (m *MockStorage) AddressExists(address string) bool {
	args := m.Called(address)
	return args.Bool(0)
}

func (m *MockStorage) SetFirstParsedBlock(address string, blockNumber uint64) {
	m.Called(address, blockNumber)
}

func (m *MockStorage) GetAllAddresses() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStorage) ReadTransactions(address string) ([]types.Transaction, error) {
	args := m.Called(address)
	return args.Get(0).([]types.Transaction), args.Error(1)
}
