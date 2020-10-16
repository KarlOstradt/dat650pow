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

// NonceHash contains a nonce and the hash
type NonceHash struct {
	Nonce int
	Hash  []byte
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
func (pow *ProofOfWork) Run(start, nRoutines int, notifyChan chan NonceHash) {
	num := big.NewInt(0)
	header := pow.setupHeader()
	for nonce := start; nonce < maxNonce; nonce += nRoutines {
		sum := sha256.Sum256(addNonce(nonce, header))
		num.SetBytes(sum[:])

		select {
		case <-notifyChan:
			return
		default:
			if num.Cmp(pow.target) == -1 {
				defer func() {
					if err := recover(); err != nil {
					}
					return
				}()

				notifyChan <- NonceHash{nonce, sum[:]}

				// close(notifyChan)

				return
			}
		}
	}
	// TODO(student)
	notifyChan <- NonceHash{0, []byte{}}
	close(notifyChan)
}

// Validate validates block's Proof-Of-Work
// This function just validates if the block header hash
// is less than the target.
func (pow *ProofOfWork) Validate() bool {
	// TODO(student)
	num := big.NewInt(0)
	sum := sha256.Sum256(addNonce(pow.block.Nonce, pow.setupHeader()))
	num.SetBytes(sum[:])

	isSmaller := num.Cmp(pow.target) == -1
	correctHash := bytes.Compare(sum[:], pow.block.Hash) == 0

	return isSmaller && correctHash
}
