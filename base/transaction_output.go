package base

import (
	"bytes"
	"fmt"
)

// TXOutput represents a transaction output
type TXOutput struct {
	Value      int    // The amount
	PubKeyHash []byte // The hash of the public key (used to "lock" the output)
}

// Lock locks the transaction to a specific address
// Only this address owns this transaction
func (out *TXOutput) Lock(address string) {
	// "Lock" the TXOutput to a specific PubKeyHash
	// based on the given address
	if !(out.PubKeyHash == nil || len(out.PubKeyHash) == 0) {
		return // Should not sign if already signed
	}
	pubKeyHash := Base58Decode([]byte(address))
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	// Create a new locked TXOutput
	pubKeyHash := Base58Decode([]byte(address))
	return &TXOutput{value, pubKeyHash[1 : len(pubKeyHash)-4]}
}

func (out TXOutput) String() string {
	return fmt.Sprintf("{%d, %x}", out.Value, out.PubKeyHash)
}
