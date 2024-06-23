package subscriber

import (
	"parser/internal/mocks"
	"parser/pkg/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSubscriber(t *testing.T) {
	// last parsed block is a reference to inform the user when we begin indexing transactions
	mockStorage := new(mocks.MockStorage)
	mockClient := new(mocks.MockClient)
	blockNumber := 100

	mockClient.On("GetBlockNumber").Return(uint64(blockNumber), nil)

	subscriber := NewSubscriber(mockStorage, mockClient)

	assert.Equal(t, uint64(blockNumber), subscriber.GetLastParsedBlock())
	mockClient.AssertExpectations(t)
}

func TestSubscriber_Subscribe(t *testing.T) {
	// when an address has not yet been subscribed, then calling subscribe should return true
	mockStorage := new(mocks.MockStorage)
	mockClient := new(mocks.MockClient)

	blockNumber := 100
	mockClient.On("GetBlockNumber").Return(uint64(blockNumber), nil)
	subscriber := NewSubscriber(mockStorage, mockClient)

	address := strings.ToLower("0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5")
	mockStorage.On("AddressExists", address).Return(false)
	mockStorage.On("SetFirstParsedBlock", address, uint64(blockNumber+1)).Return(nil)

	subscribed := subscriber.Subscribe(address)
	assert.True(t, subscribed)
	mockStorage.AssertExpectations(t)
}

func TestSubscriber_GetSubscribedAddresses(t *testing.T) {
	// user should be able to retrieve all subscribed addresses
	mockStorage := new(mocks.MockStorage)
	mockClient := new(mocks.MockClient)

	blockNumber := 100
	mockClient.On("GetBlockNumber").Return(uint64(blockNumber), nil)
	expectedAddresses := []string{"0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"}
	mockStorage.On("GetAllAddresses").Return(expectedAddresses, nil)

	subscriber := NewSubscriber(mockStorage, mockClient)

	addresses, err := subscriber.GetSubscribedAddresses()
	assert.NoError(t, err)
	assert.Equal(t, expectedAddresses, addresses)
	mockStorage.AssertExpectations(t)
}

func TestSubscriber_GetTransactions(t *testing.T) {
	// retrieve transactions from storage
	mockStorage := new(mocks.MockStorage)
	mockClient := new(mocks.MockClient)

	blockNumber := 100
	mockClient.On("GetBlockNumber").Return(uint64(blockNumber), nil)

	var tx types.Transaction = map[string]interface{}{
		"hash":      "0xb5c8bd9430b6cc87a0e2fe110ece6bf527fa4f170a4bc8cd032f768fc5219838",
		"blockHash": "0x1",
	}
	expectedTransactions := []types.Transaction{tx}
	mockStorage.On("ReadTransactions", "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5").Return(expectedTransactions, nil)

	subscriber := NewSubscriber(mockStorage, mockClient)

	transactions := subscriber.GetTransactions("0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5")
	assert.Equal(t, expectedTransactions, transactions)
	mockStorage.AssertExpectations(t)
}
