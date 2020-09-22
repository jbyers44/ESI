package trie

import (
	"crypto/sha256"
)

// Leaf is the data structure for terminating leaf nodes that only contain a value and a hash of that value
type Leaf struct {
	value []byte
	hash  []byte
}

// NewLeaf is the default contructor for a leaf
func NewLeaf(value []byte) *Leaf {
	h := sha256.New()
	h.Write(value)
	return &Leaf{value, h.Sum(nil)}
}
