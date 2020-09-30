package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Fixed block timestamp
const BlockTime int64 = 1563897484

func TestGenesisBlock(t *testing.T) {
	// Genesis block
	gb := NewGenesisBlock(testTransactions["tx0"])
	assert.NotNil(t, gb)

	assert.Equal(t, []byte{}, gb.PrevBlockHash, "Genesis block shouldn't has PrevBlockHash")

	assert.Equal(t, testTransactions["tx0"].ID, gb.Transactions[0].ID, "Genesis block should contains the coinbase transaction")
}

func TestBlockHashTransactions(t *testing.T) {
	// Merkle root of txs
	merkleRootTxsHash := Hex2Bytes("153b097cb029051834d285f69901237812f13d6dc166d684778096422f331677")

	b := &Block{0, []*Transaction{testTransactions["tx1"]}, []byte{}, []byte{}, 0}

	assert.Equalf(t, merkleRootTxsHash, b.HashTransactions(), "The block hash %x isn't equal to %x", b.HashTransactions(), merkleRootTxsHash)
}

func TestMine(t *testing.T) {
	genesisBlock := testBlockchainData["block0"]

	b := &Block{BlockTime, []*Transaction{testTransactions["tx1"]}, genesisBlock.Hash, []byte{}, 0}
	b.Mine()

	assert.Equalf(t, testBlockchainData["block1"].Hash, b.Hash, "The block hash %x isn't equal to %x", b.Hash, testBlockchainData["block1"].Hash)
	assert.Equalf(t, testBlockchainData["block1"].Nonce, b.Nonce, "The block nonce %d isn't equal to %d", b.Nonce, testBlockchainData["block1"].Nonce)
}

func TestNewBlock(t *testing.T) {
	genesisBlock := testBlockchainData["block0"]
	b := NewBlock([]*Transaction{testTransactions["tx1"]}, genesisBlock.Hash)

	assert.NotEqual(t, []byte{}, b.Hash, "The block hash should have a valid value")

	assert.Equal(t, genesisBlock.Hash, b.PrevBlockHash, "Previous block of the current should be the genesis block")
}
