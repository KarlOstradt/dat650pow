package base

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// TARGETBITS define the mining difficulty
const TARGETBITS = 16

// ProofOfWork represents a block mined with a target difficulty
type ProofOfWork struct {
	block  *Block
	target *big.Int // 2 ** (256 - targetbits)
}

// NewProofOfWork builds a ProofOfWork
func NewProofOfWork(block *Block) *ProofOfWork {
	// TODO(student)
	x := *big.NewInt(0)
	x.SetBit(&x, 256-TARGETBITS, 1)
	return &ProofOfWork{block, &x}
}

// setupHeader prepare the header of the block
func (pow *ProofOfWork) setupHeader() []byte {
	// TODO(student)
	var header []byte
	header = append(header, pow.block.PrevBlockHash...)
	header = append(header, pow.block.HashTransactions()...)
	header = append(header, IntToHex(pow.block.Timestamp)...)
	header = append(header, IntToHex(TARGETBITS)...)

	return header
}

// addNonce adds a nonce to the header
func addNonce(nonce int, header []byte) []byte {
	// TODO(student)
	return append(header, IntToHex(int64(nonce))...)
}

// Run performs the proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	num := big.NewInt(0)
	header := pow.setupHeader()
	for nonce := 0; nonce < maxNonce; nonce++ {

		sum := sha256.Sum256(addNonce(nonce, header))
		num.SetBytes(sum[:])
		if num.Cmp(pow.target) == -1 {
			return nonce, sum[:]
		}
		// if bytes.Compare(sum[:], pow.target.Bytes()) < 0 {
		// 	return nonce, sum[:]
		// }
	}
	// TODO(student)
	return 0, []byte{}
}

// Validate validates block's Proof-Of-Work
// This function just validates if the block header hash
// is less than the target.
func (pow *ProofOfWork) Validate() bool {
	// TODO(student)
	isSmaller := bytes.Compare(pow.block.Hash, pow.target.Bytes()) < 0

	sum := sha256.Sum256(addNonce(pow.block.Nonce, pow.setupHeader()))
	validNonce := bytes.Compare(sum[:], pow.target.Bytes()) < 0
	return isSmaller && validNonce
}
