package base

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Transaction represents a Bitcoin transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)

	enc.Encode(Transaction{tx.ID, tx.Vin, tx.Vout})

	return data.Bytes()
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	t := *tx
	t.ID = []byte{}
	sum := sha256.Sum256(t.Serialize())
	return sum[:]
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       OutIdx:    %d", input.OutIdx))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey: %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	hasOneInput := len(tx.Vin) == 1
	if !hasOneInput {
		return false
	}
	input := tx.Vin[0]
	output := tx.Vout[0]

	hasEmptyID := len(input.Txid) == 0
	isFirst := input.OutIdx == -1
	hasNoSignature := input.Signature == nil

	hasCorrectReward := output.Value == BlockReward

	return hasOneInput && hasEmptyID && isFirst && hasCorrectReward && hasNoSignature
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = "Reward to " + to
	}
	tx := &Transaction{
		Vin: []TXInput{
			{Txid: []byte{}, OutIdx: -1, Signature: nil, PubKey: []byte(data)},
		},
		Vout: []TXOutput{
			*NewTXOutput(BlockReward, to),
		},
	}
	tx.ID = tx.Hash()
	return tx
}

// NewUTXOTransaction creates a new UTXO transaction
func NewUTXOTransaction(wallet *Wallet, to string, amount int, utxos UTXOSet, bc *Blockchain) (*Transaction, error) {
	hashedPubKey := GetPubKeyHashFromAddress(wallet.GetStringAddress())
	n, spendableOutputs := utxos.FindSpendableOutputs(hashedPubKey, amount)
	if n < amount {
		return nil, errors.New("Not enough funds")
	}

	var inputs []TXInput
	var outputs []TXOutput
	prevTXs := make(map[string]*Transaction)
	// Create inputs
	for txID, indexes := range spendableOutputs {
		for _, i := range indexes {
			inn := TXInput{Txid: Hex2Bytes(txID), OutIdx: i, PubKey: wallet.PublicKey, Signature: nil}
			inputs = append(inputs, inn)

		}
		prevTX, _ := bc.FindTransaction(Hex2Bytes(txID))
		prevTXs[txID] = prevTX
	}

	//Create outputs
	output := *NewTXOutput(amount, to)
	outputs = append(outputs, output)
	if n-amount > 0 {
		change := *NewTXOutput(n-amount, wallet.GetStringAddress())
		outputs = append(outputs, change)
	}

	// Create the transaction
	tx := &Transaction{Vin: inputs, Vout: outputs}
	tx.ID = tx.Hash()
	tx.Sign(wallet.PrivateKey, prevTXs)
	return tx, nil
}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	for _, input := range tx.Vin {
		inputs = append(inputs, TXInput{Txid: input.Txid, OutIdx: input.OutIdx})
	}
	return Transaction{ID: tx.ID, Vin: inputs, Vout: tx.Vout}
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	// 1) coinbase transactions are not signed.
	// 2) Throw a Panic in case of any prevTXs (used inputs) didn't exists
	// Take a look on the tests to see the expected error message
	// 3) Create a copy of the transaction to be signed
	// 4) Sign all the previous TXInputs of the transaction tx using the
	// copy as the payload (serialized) to be signed in the ecdsa.Sig
	// (https://golang.org/pkg/crypto/ecdsa/#Sign)
	// Make sure that each input of the copy to be signed
	// have the correct PubKeyHash of each output in the prevTXs
	// Store the signature as a concatenation of R and S fields
	if tx.IsCoinbase() {
		return
	}

	for _, in := range tx.Vin {
		key := hex.EncodeToString(in.Txid)
		if _, ok := prevTXs[key]; !ok {
			panic("Current input transaction isn't listed in previous transactions")
		} else {
			// Do more checks?
		}
	}

	txCopy := tx.TrimmedCopy()
	for i := range tx.Vin {
		r, s, _ := ecdsa.Sign(rand.Reader, &privKey, txCopy.Serialize())
		sig := append(r.Bytes(), s.Bytes()...)
		tx.Vin[i].Signature = sig
	}
}

// Verify verifies signatures of Transaction inputs
func (tx Transaction) Verify(prevTXs map[string]*Transaction) bool {
	// 1) coinbase transactions are not signed.
	// 2) Throw a Panic in case of any prevTXs (used inputs) didn't exists
	// Take a look on the tests to see the expected error message
	// 3) Create the same copy of the transaction that was signed
	// and get the curve used for sign: P256
	// 4) Doing the opposite operation of the signing, perform the
	// verification of the signature, by recovering the R and S byte fields
	// of the Signature and the X and Y fields of the PubKey from
	// the inputs of tx. Verify the signature of each input using the
	// ecdsa.Verify function (https://golang.org/pkg/crypto/ecdsa/#Verify)
	// Note that to use this function you need to reconstruct the
	// ecdsa.PublicKey. Also notice that the ecdsa.Verify function receive
	// a byte array, you the transaction copy need to be serialized.
	// return true if all inputs have valid signature,
	// and false if any of them have an invalid signature.

	if tx.IsCoinbase() {
		return true
	}

	// Check that all inputs are referencing another transaction
	for _, in := range tx.Vin {
		key := hex.EncodeToString(in.Txid)
		if _, ok := prevTXs[key]; !ok {
			panic("Current input transaction isn't listed in previous transactions")
		} else {
			// Do more checks?
		}
	}

	// txCopy := tx.TrimmedCopy() // Why do we need the copy?
	for _, input := range tx.Vin {
		var pubKey ecdsa.PublicKey
		pubKey.Curve = elliptic.P256()

		x := big.NewInt(0)
		xBytes := input.PubKey[0 : len(input.PubKey)/2]
		x.SetBytes(xBytes)
		pubKey.X = x

		y := big.NewInt(0)
		yBytes := input.PubKey[len(input.PubKey)/2:]
		y.SetBytes(yBytes)
		pubKey.Y = y

		sig := input.Signature
		r := big.NewInt(0)
		r.SetBytes(sig[0 : len(sig)/2])
		s := big.NewInt(0)
		s.SetBytes(sig[len(sig)/2:])

		if !ecdsa.Verify(&pubKey, tx.Serialize(), r, s) {
			return false
		}
	}
	return true
}
