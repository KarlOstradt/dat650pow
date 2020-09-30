package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func newMockHeader() []byte {
	return bytes.Join(
		[][]byte{
			[]byte{},
			Hex2Bytes("fdfa9ad1db072757d55c11ba05aecae0bbd99e29b8dc2a869a68ebeb1ca09147"),
			IntToHex(BlockTime),
			IntToHex(TARGETBITS),
		},
		[]byte{},
	)
}

// TARGETBITS == 8 => target difficulty of 2^248
// Hexadecimal: 100000000000000000000000000000000000000000000000000000000000000
// Big Int: 452312848583266388373324160190187140051835877600158453279131187530910662656
var testTargetDifficulty, _ = new(big.Int).SetString("452312848583266388373324160190187140051835877600158453279131187530910662656", 10)

func TestNewProofOfWork(t *testing.T) {
	b := &Block{BlockTime, []*Transaction{testTransactions["tx0"]}, []byte{}, []byte{}, 0}

	pow := NewProofOfWork(b)

	assert.Equal(t, b, pow.block)
	assert.Equal(t, testTargetDifficulty, pow.target)
}

func TestSetupHeader(t *testing.T) {
	b := &Block{BlockTime, []*Transaction{testTransactions["tx0"]}, []byte{}, []byte{}, 0}

	pow := &ProofOfWork{b, testTargetDifficulty}

	header := pow.setupHeader()

	expectedHeader := newMockHeader()
	assert.Equalf(t, expectedHeader, header, "The current block header: %x isn't equal to the expected %x\n", header, expectedHeader)
}

func TestAddNonce(t *testing.T) {
	header := newMockHeader()
	expectedHeader := bytes.Join(
		[][]byte{
			header,
			IntToHex(33),
		},
		[]byte{},
	)
	assert.Equal(t, expectedHeader, addNonce(33, header))
}

func TestRun(t *testing.T) {
	for i, block := range testBlockchainData {
		b := &Block{BlockTime, block.Transactions, block.PrevBlockHash, []byte{}, 0}
		pow := &ProofOfWork{b, testTargetDifficulty}
		nonce, hash := pow.Run()
		assert.Equal(t, testBlockchainData[i].Nonce, nonce)
		assert.Equal(t, testBlockchainData[i].Hash, hash)
	}
}
func TestValidatePoW(t *testing.T) {
	for _, block := range testBlockchainData {
		pow := &ProofOfWork{block, testTargetDifficulty}
		assert.True(t, pow.Validate())
	}
}
