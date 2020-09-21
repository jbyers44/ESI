package trei

import (
	"crypto/sha256"
)

// Leaf is the data structure for terminating leaf nodes that only contain a value and a hash of that value
type Leaf struct {
	value []byte
	hash  [32]byte
}

// NewLeaf is the default contructor for a leaf
func NewLeaf(value []byte) *Leaf {
	hash := sha256.Sum256(value)
	return &Leaf{value, hash}
}
