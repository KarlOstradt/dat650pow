package base

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	RootNode *Node
	Leafs    []*Node
}

// Node represents a Merkle tree node
type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	Hash   []byte
}

const (
	leftNode = iota
	rightNode
)

// MerkleProof represents way to prove element inclusion on the merkle tree
type MerkleProof struct {
	proof [][]byte
	index []int64
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
func NewMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0 {
		panic("No merkle tree nodes")
	}

	var leafs []*Node
	for _, leaf := range data {
		leafs = append(leafs, NewMerkleNode(nil, nil, leaf))
	}

	root := constructLayers(leafs)[0]

	return &MerkleTree{RootNode: root, Leafs: leafs}
}

func constructLayers(nodes []*Node) []*Node {
	var layer []*Node

	if len(nodes) == 1 {
		return nodes
	}

	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		var right *Node
		if i+1 < len(nodes) {
			right = nodes[i+1]
		}
		var data []byte
		if right == nil {
			data = append(left.Hash, left.Hash...)
		} else {
			data = append(left.Hash, right.Hash...)

		}
		layer = append(layer, NewMerkleNode(left, right, data))
	}
	return constructLayers(layer)
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *Node, data []byte) *Node {
	var sum [32]byte
	if left == nil && right == nil {
		sum = sha256.Sum256(data)
		return &Node{Left: left, Right: right, Hash: sum[:]}
	}

	if right == nil {
		sum = sha256.Sum256(append(left.Hash, left.Hash...))
	} else {
		sum = sha256.Sum256(append(left.Hash, right.Hash...))
	}

	parent := &Node{Left: left, Right: right, Hash: sum[:]}
	if left != nil {
		left.Parent = parent
	}
	if right != nil {
		right.Parent = parent
	}
	return parent
}

// MerkleRootHash return the hash of the merkle root
func (mt *MerkleTree) MerkleRootHash() []byte {
	return mt.RootNode.Hash
}

// MakeMerkleProof returns a list of hashes and indexes required to
// reconstruct the merkle path of a given hash
//
// @param hash represents the hashed data (e.g. transaction ID) stored on
// the leaf node
// @return the merkle proof (list of intermediate hashes), a list of indexes
// indicating the node location in relation with its parent (using the
// constants: leftNode or rightNode), and a possible error.
func (mt *MerkleTree) MakeMerkleProof(hash []byte) ([][]byte, []int64, error) {
	var node *Node
	for _, leaf := range mt.Leafs {
		if bytes.Equal(leaf.Hash, hash) {
			node = leaf
			break
		}
	}

	if node == nil {
		return [][]byte{}, []int64{}, fmt.Errorf("Node %x not found", hash)
	}

	path := [][]byte{}
	index := []int64{}

	for node.Parent != nil {
		// Left
		if bytes.Equal(node.Parent.Left.Hash, node.Hash) {
			if node.Parent.Right == nil {
				path = append(path, node.Parent.Left.Hash)
			} else {
				path = append(path, node.Parent.Right.Hash)
			}
			index = append(index, rightNode)
		} else {
			path = append(path, node.Parent.Left.Hash)
			index = append(index, leftNode)
		}
		node = node.Parent
	}

	return path, index, nil
}

// VerifyProof verifies that the correct root hash can be retrieved by
// recreating the merkle path for the given hash and merkle proof.
//
// @param rootHash is the hash of the current root of the merkle tree
// @param hash represents the hash of the data (e.g. transaction ID)
// to be verified
// @param mProof is the merkle proof that contains the list of intermediate
// hashes and their location on the tree required to reconstruct
// the merkle path.
func VerifyProof(rootHash []byte, hash []byte, mProof MerkleProof) bool {
	if len(mProof.proof) != len(mProof.index) {
		return false
	}
	if len(mProof.proof) == 0 && bytes.Equal(rootHash, hash) {
		return true
	}

	tempH := hash
	for i, h := range mProof.proof {
		var sum [32]byte
		if mProof.index[i] == 0 { // h is leftNode
			sum = sha256.Sum256(append(h, tempH...))
		} else { // h is rightNode
			sum = sha256.Sum256(append(tempH, h...))
		}

		tempH = sum[:]
	}

	return bytes.Equal(rootHash, tempH)
}
