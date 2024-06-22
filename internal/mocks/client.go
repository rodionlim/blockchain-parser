package mocks

import (
	"trustwallet/pkg/types"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetBlockNumber() (uint64, error) {
	args := m.Called()
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockClient) GetBlockByNumber(blockNumber uint64) (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockClient) GetTransactionsBySubscriptionsAndBlock(subscriptions map[string]*types.Tracker, blockNumber uint64) (map[string][]types.Transaction, error) {
	args := m.Called(subscriptions, blockNumber)
	return args.Get(0).(map[string][]types.Transaction), args.Error(1)
}
