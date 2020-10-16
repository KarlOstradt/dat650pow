package base

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Block keeps block information
type Block struct {
	Timestamp     int64          // the block creation timestamp
	Transactions  []*Transaction // The block transactions
	PrevBlockHash []byte         // the hash of the previous block
	Hash          []byte         // the hash of the block
	Nonce         int            // the nonce of the block
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Transactions: transactions, PrevBlockHash: prevBlockHash}
	block.Mine(nil, 0, 1) // will set hash and nonce
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// Mine calculates and sets the block hash and nonce.
func (b *Block) Mine(stopChan chan bool, id, nSlaves int) {
	// TODO(student)
	pow := NewProofOfWork(b)
	notifyChan := make(chan NonceHash)
	for start := id; start < nRoutines*nSlaves; start += nSlaves {
		go pow.Run(start, nRoutines*nSlaves, notifyChan)
	}
	select {
	case <-stopChan:
		notifyChan <- NonceHash{}
		// close(notifyChan)
	case nh := <-notifyChan:
		b.Nonce = nh.Nonce
		b.Hash = nh.Hash
	}
	close(notifyChan)
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var data [][]byte
	for _, tx := range b.Transactions {
		data = append(data, tx.Serialize())
	}
	merkleTree := NewMerkleTree(data)
	return merkleTree.RootNode.Hash
}

// FindTransaction finds a transaction by its ID
func (b *Block) FindTransaction(ID []byte) (*Transaction, error) {
	for _, tx := range b.Transactions {
		if bytes.Equal(tx.ID, ID) {
			return tx, nil
		}
	}
	return nil, errors.New("Transaction not found")
}

func (b *Block) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("============ Block %x ============", b.Hash))
	lines = append(lines, fmt.Sprintf("Prev. hash: %x", b.PrevBlockHash))
	lines = append(lines, fmt.Sprintf("Timestamp: %v\n", time.Unix(b.Timestamp, 0)))
	for _, tx := range b.Transactions {
		lines = append(lines, fmt.Sprintf("%v\n", tx))
	}
	return strings.Join(lines, "\n")
}
