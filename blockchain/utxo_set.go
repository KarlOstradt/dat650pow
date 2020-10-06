package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

// UTXOSet represents a set of UTXO as an in-memory cache
// The key of the map is the transaction ID
// (encoded as string) that contains these outputs
type UTXOSet map[string][]TXOutput

// FindSpendableOutputs finds and returns unspent outputs in the UTXO Set
// to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	// Modify your function to use IsLockedWithKey instead of CanBeUnlockedWith
	spendable := make(map[string][]int)
	sum := 0
	for txID, outputs := range u {
		var out []int
		for i, output := range outputs {

			if (&output).IsLockedWithKey(pubKeyHash) {
				if sum < amount { // Only append enough outputs
					out = append(out, i)
					sum += output.Value
				} else {
					break
				}
			}
		}
		spendable[txID] = out
	}

	return sum, spendable
}

// FindUTXO finds UTXO in the UTXO Set for a given unlockingData key (e.g., address)
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	// Modify your function to use IsLockedWithKey instead of CanBeUnlockedWith
	var UTXO []TXOutput
	for _, outputs := range u {
		for _, output := range outputs {
			if (&output).IsLockedWithKey(pubKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}

// CountUTXOs returns the number of transactions outputs in the UTXO set
func (u UTXOSet) CountUTXOs() int {
	count := 0
	for _, outputs := range u {
		count += len(outputs)
	}
	return count
}

// Update updates the UTXO Set with the new set of transactions
func (u UTXOSet) Update(transactions []*Transaction) {
	for _, tx := range transactions {

		// Find and delete outputs that are linked to inputs
		for _, input := range tx.Vin {
			key := hex.EncodeToString(input.Txid)
			if _, ok := u[key]; ok {
				var newOutputs []TXOutput
				for i, output := range u[key] {
					if i != input.OutIdx {
						newOutputs = append(newOutputs, output)
					}
				}

				if len(newOutputs) == 0 {
					delete(u, key)
				} else {
					u[key] = newOutputs
				}

			}
		}

		// Add all output of this tx as UTXO for now
		for _, output := range tx.Vout {
			key := hex.EncodeToString(tx.ID)
			u[key] = append(u[key], output)
		}
	}
}

// Equal compares two UTXOSet
func (u UTXOSet) Equal(utxos UTXOSet) bool {
	if len(u) != len(utxos) {
		return false
	}

	for txid, outputs := range u {
		o, ok := utxos[txid]

		if !ok || len(outputs) != len(o) {
			return false
		}

		for i, out := range outputs {
			if out.Value != o[i].Value || !bytes.Equal(out.PubKeyHash, o[i].PubKeyHash) {
				return false
			}
		}
	}

	return true
}

func (u UTXOSet) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- UTXO SET:"))
	for txid, outputs := range u {
		lines = append(lines, fmt.Sprintf("     TxID: %s", txid))
		for i, out := range outputs {
			lines = append(lines, fmt.Sprintf("           Output %d: %v", i, out))
		}
	}

	return strings.Join(lines, "\n")
}
