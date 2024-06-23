package storage

import (
	"parser/pkg/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_ReadTransactions(t *testing.T) {
	storage := NewMemoryStorage()

	address := "0x123"
	lowercaseAddress := strings.ToLower(address)

	var tx1 types.Transaction = map[string]interface{}{"hash": "0xabc"}
	var tx2 types.Transaction = map[string]interface{}{"hash": "0xdef"}
	var tx3 types.Transaction = map[string]interface{}{"hash": "0xyz"}
	var tx4 types.Transaction = map[string]interface{}{"hash": "0xyz2"}

	transactions := []types.Transaction{tx1, tx2}

	// Set initial state in memory storage
	storage.SetFirstParsedBlock(address, 0)
	storage.WriteTransactions(map[string][]types.Transaction{
		address: transactions,
	})

	// Read transactions
	readTxns, err := storage.ReadTransactions(address)
	require.NoError(t, err, "unexpected error while reading transactions")
	assert.Equal(t, transactions, readTxns, "transactions do not match")

	// Verify that the offset is updated correctly
	expectedOffset := len(transactions)
	actualOffset := storage.trackerData[lowercaseAddress].Offset
	assert.Equal(t, expectedOffset, actualOffset, "offset not updated correctly")

	// Add more transactions
	moreTransactions := []types.Transaction{tx3, tx4}
	storage.WriteTransactions(map[string][]types.Transaction{
		lowercaseAddress: moreTransactions,
	})

	// // Read new transactions
	readTxns, err = storage.ReadTransactions(address)
	require.NoError(t, err, "unexpected error while reading transactions")
	assert.Equal(t, moreTransactions, readTxns, "new transactions do not match")

	// Verify that the offset is updated correctly again
	expectedOffset = len(transactions) + len(moreTransactions)
	actualOffset = storage.trackerData[lowercaseAddress].Offset
	assert.Equal(t, expectedOffset, actualOffset, "offset not updated correctly after adding more transactions")
}

func TestMemoryStorage_ReadTransactions_InvalidAddress(t *testing.T) {
	storage := NewMemoryStorage()

	address := "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"

	_, err := storage.ReadTransactions(address)
	require.Error(t, err, "expected error for non-existent address")
	assert.Contains(t, err.Error(), "invalid address [0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5] specified")
}

func TestMemoryStorage_ReadTransactions_NoTransactions(t *testing.T) {
	storage := NewMemoryStorage()

	address := "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"

	// Set initial state in memory storage
	storage.SetFirstParsedBlock(address, 0)

	// Read transactions (no transactions should be present)
	readTxns, err := storage.ReadTransactions(address)
	require.NoError(t, err, "unexpected error while reading transactions")
	assert.Nil(t, readTxns, "expected nil when no transactions present")
}
