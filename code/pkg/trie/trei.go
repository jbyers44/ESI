package trie

import (
	"bytes"
	"crypto/sha256"
)

// MerklePatriciaTrie is effectively a wrapper for a pointer to the root node
type MerklePatriciaTrie struct {
	root *Node
}

// NewMerklePatriciaTrie creates a new trie with a root node with two empty string labels
func NewMerklePatriciaTrie() *MerklePatriciaTrie {
	return &MerklePatriciaTrie{nil}
}

// Insert a node into the trie
func (trie *MerklePatriciaTrie) Insert(value []byte) {
	if trie.root == nil {
		leaf := NewLeaf(value)
		trie.root = &Node{[32]byte{}, []byte{}, leaf, []byte{}, nil}
	} else {
		insert(trie.root, value, []byte(""))
	}
}

func makeNode(value []byte, label []byte, suffixEndIndex int) *Node {
	newNodeLeftLabel := value[len(label):suffixEndIndex][:]
	newLeftLeaf := &Leaf{value, [32]byte{}}

	newNodeRightLabel := []byte("")
	newRightLeaf := &Leaf{label, [32]byte{}}

	suffixHash := sha256.Sum256(value[0:len(label)][:])
	newLeftLabelHash := sha256.Sum256(newNodeLeftLabel)
	// newNodeLeftLabel := value[len(root.leftLabel):suffixEnd][:]
	// newLeftLeaf := &Leaf{value, [32]byte{}}
	// newNodeRightLabel := []byte("")
	// newRightLeaf := &Leaf{root.leftLabel, [32]byte{}}
	hash := sha256.Sum256(append(suffixHash[:], newLeftLabelHash[:]...))

	newNode := &Node{hash, newNodeLeftLabel, newLeftLeaf, newNodeRightLabel, newRightLeaf}
	return newNode
}

// Private recursive func that does actual insertion logic
func insert(root *Node, value []byte, prefix []byte) {
	if bytes.HasPrefix(value, root.leftLabel) {
		suffixEnd := len(root.leftLabel) + len(value) - 2
		prefix := value[0:suffixEnd][:]
		if _, isLeaf := root.left.(Leaf); isLeaf {
			// newLeftLabelHash := sha256.Sum256(newNodeLeftLabel)
			// newNodeLeftLabel := value[len(root.leftLabel):suffixEnd][:]
			// newLeftLeaf := &Leaf{value, [32]byte{}}
			// newNodeRightLabel := []byte("")
			// newRightLeaf := &Leaf{root.leftLabel, [32]byte{}}
			// hash := sha256.Sum256(append(suffixHash, newLeftLabelHash...))
			newNode := makeNode(value, root.leftLabel, suffixEnd)
			root.SetLeft(newNode)
		} else if nextNode, isNode := root.left.(Node); isNode {
			nextPrefix := append(prefix, root.leftLabel...)
			insert(&nextNode, value, nextPrefix)
		}
	} else if bytes.HasPrefix(value, root.rightLabel) {
		suffixEnd := len(root.rightLabel) + len(value) - 2
		prefix := value[0:suffixEnd][:]
		if _, isLeaf := root.right.(Leaf); isLeaf {
			newNode := makeNode(value, root.rightLabel, suffixEnd)
			root.SetRight(newNode)
		} else if nextNode, isNode := root.right.(Node); isNode {
			nextPrefix := append(prefix, root.rightLabel...)
			insert(&nextNode, value, nextPrefix)
		}
	}

	print("llabel: " + string(root.leftLabel))
	print("rlabel: " + string(root.rightLabel))

}

// NewTrie makes a new trie
func (trie *MerklePatriciaTrie) NewTrie(values [][]byte) {
	for i := range values {
		trie.Insert(values[i])
	}
}

// GetRoot gets the root
func (trie *MerklePatriciaTrie) GetRoot() (root *Node) {
	return trie.root
}
