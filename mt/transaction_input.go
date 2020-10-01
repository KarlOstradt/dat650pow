package main

import (
	"bytes"
	"crypto/sha256"
)

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte // The ID (i.e. Hash) of the transaction.
	OutIdx    int    // The index of the output
	Signature []byte // The signature of this input.
	PubKey    []byte // The raw public key (not hashed)
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	sum := sha256.Sum256(in.PubKey)
	return bytes.Equal(sum[:], pubKeyHash)
}
