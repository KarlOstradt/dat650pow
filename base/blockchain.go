package base

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	blocks []*Block
}

// CreateBlockchain creates a new blockchain with genesis Block
func CreateBlockchain(address string) *Blockchain {
	inn := TXInput{Txid: []byte{}, OutIdx: -1, Signature: nil, PubKey: []byte(GenesisCoinbaseData)}
	out := *NewTXOutput(BlockReward, address)
	tx := Transaction{Vin: []TXInput{inn}, Vout: []TXOutput{out}}
	tx.ID = tx.Hash()
	genesisBlock := NewGenesisBlock(&tx)
	blockchain := Blockchain{blocks: []*Block{genesisBlock}}

	return &blockchain
}

// NewBlockchain creates a Blockchain
func NewBlockchain(address string) *Blockchain {
	return CreateBlockchain(address)
}

// AddBlock saves the block into the blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) *Block {
	current := bc.CurrentBlock()
	block := NewBlock(transactions, current.Hash)
	bc.blocks = append(bc.blocks, block)
	return block
}

// GetGenesisBlock returns the Genesis Block
func (bc Blockchain) GetGenesisBlock() *Block {
	return bc.blocks[0]
}

// CurrentBlock returns the last block
func (bc Blockchain) CurrentBlock() *Block {
	return bc.blocks[len(bc.blocks)-1]
}

// GetBlock returns the block of a given hash
func (bc Blockchain) GetBlock(hash []byte) (*Block, error) {
	for i := len(bc.blocks) - 1; i >= 0; i-- {
		if bytes.Equal(bc.blocks[i].Hash, hash) {
			return bc.blocks[i], nil
		}
	}

	return nil, errors.New("no blocks has the given hash")
}

// ValidateBlock validates the a block after mining or
// before adding it to the blockchain
func (bc *Blockchain) ValidateBlock(block *Block) bool {
	// TODO(student)
	// check if and only if the first tx is coinbase
	// validates block's Proof-Of-Work
	if !block.Transactions[0].IsCoinbase() {
		return false
	}
	if bytes.Compare(block.PrevBlockHash, bc.CurrentBlock().Hash) != 0 {
		return false
	}
	pow := NewProofOfWork(block)

	return pow.Validate()
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) (*Block, error) {
	var validTxs []*Transaction
	for _, tx := range transactions {
		isValid := bc.VerifyTransaction(tx)
		if isValid {
			validTxs = append(validTxs, tx)
		}
	}

	if len(validTxs) == 0 {
		return nil, errors.New("there are no valid transactions to be mined")
	}

	return bc.AddBlock(validTxs), nil
}

// FindTransaction finds a transaction by its ID in the whole blockchain
func (bc Blockchain) FindTransaction(ID []byte) (*Transaction, error) {
	for _, block := range bc.blocks {
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return tx, nil
			}
		}
	}

	return nil, errors.New("Transaction not found in any block")
}

// FindUTXOSet finds and returns all unspent transaction outputs
func (bc Blockchain) FindUTXOSet() UTXOSet {
	UTXO := make(UTXOSet)
	for _, block := range bc.blocks {
		UTXO.Update(block.Transactions)
	}

	return UTXO
}

// GetInputTXsOf returns a map index by the ID,
// of all transactions used as inputs in the given transaction
func (bc *Blockchain) GetInputTXsOf(tx *Transaction) (map[string]*Transaction, error) {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	// Use bc.FindTransaction to search over all transactions
	// in the blockchain and if the referred input into tx exists,
	// if so, get the transaction of this input and add it
	// to a map, where the key is the id of the transaction found
	// and the value is the pointer to transaction itself.
	// To use the id as key in the map, convert it to string
	// using the function: hex.EncodeToString
	// https://golang.org/pkg/encoding/hex/#EncodeToString

	prevTXs := make(map[string]*Transaction)

	for _, input := range tx.Vin {
		key := hex.EncodeToString(input.Txid)
		if _, ok := prevTXs[key]; !ok {
			prevTX, err := bc.FindTransaction(input.Txid)
			if err != nil {
				// panic(err.Error())
				return nil, err
			}
			prevTXs[key] = prevTX
		}
	}

	return prevTXs, nil
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	// Get the previous transactions referred in the input of tx
	// and call Sign for tx.

	prevTXs, err := bc.GetInputTXsOf(tx)
	if err != nil {
		panic(err.Error())
	}
	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	// Modify the function to get the inputs referred in tx
	// and return false in case of some error (i.e. not found the input).
	// Then call Verify for tx passing those inputs as parameter and return the result.
	// Remember that coinbase transaction doesn't have input.

	if tx.IsCoinbase() {
		return true
	}
	prevTXs, err := bc.GetInputTXsOf(tx)
	if err != nil {
		return false
	}
	return tx.Verify(prevTXs)
}

func (bc Blockchain) String() string {
	var lines []string
	for _, block := range bc.blocks {
		lines = append(lines, fmt.Sprintf("%v", block))
	}
	return strings.Join(lines, "\n")
}
