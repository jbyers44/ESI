package trie

import (
	"bytes"
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
		dummy := NewLeaf([]byte(""))
		trie.root = &Node{[32]byte{}, value, leaf, []byte(""), dummy}
	} else {
		insert(trie.root, value, 0)
	}
}

func greatestCommonPrefix(prefix []byte, value []byte) int {
	i := 0
	for i < len(prefix) && i < len(value) && prefix[i] == value[i] {
		i++
	}
	return i
}

func insertNode(node interface{}, label []byte, value []byte, gcp int, prefix int) *Node {
	labelRemainder := label[gcp:]
	valueRemainder := value[prefix+gcp:]
	newLeaf := NewLeaf(value)
	return &Node{[32]byte{}, labelRemainder, node, valueRemainder, newLeaf}
}

// Private recursive func that does actual insertion logic
func insert(root *Node, value []byte, prefix int) {
	//  Greatest common prefix
	gcp := 0

	// Handle left pointer
	gcp = greatestCommonPrefix(value[prefix:], root.leftLabel)

	if gcp > 0 {
		switch v := root.left.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on right
			if bytes.Compare(v.value, value) == 0 {
				return
			}
			root.left = insertNode(v, root.leftLabel, value, gcp, prefix)
			root.leftLabel = root.leftLabel[:gcp]
			return

		case *Node:
			if gcp == len(root.leftLabel) { // If the prefix is a full match, traverse down to the next node
				insert(v, value, prefix+gcp)
				return
			} // The prefix isn't a full match, create a new intermediary node with a new leaf value on right, old tree on left
			root.left = insertNode(v, root.leftLabel, value, gcp, prefix)
			root.leftLabel = root.leftLabel[:gcp]
			return
		}
	}

	// Handle right pointer
	gcp = greatestCommonPrefix(value[prefix:], root.rightLabel)

	if gcp > 0 {

		switch v := root.right.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on right
			if bytes.Compare(v.value, value) == 0 {
				return
			}
			root.right = insertNode(v, root.rightLabel, value, gcp, prefix)
			root.rightLabel = root.rightLabel[:gcp]
			return

		case *Node:
			if gcp == len(root.rightLabel) { // If the prefix is a full match, traverse down to the next node
				insert(v, value, prefix+gcp)
				return
			} // The prefix isn't a full match, create a new intermediary node with a new leaf value on right, old tree on left
			root.right = insertNode(v, root.rightLabel, value, gcp, prefix)
			root.rightLabel = root.rightLabel[:gcp]
			return
		}
	} else { // Final case, shares no prefix with either node

		switch v := root.right.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on right
			if bytes.Compare(v.value, []byte("")) == 0 {
				root.right = NewLeaf(value)
				root.rightLabel = value
				return
			}
			goto ExtensionNode

		case *Node:
			if bytes.Compare(root.rightLabel, []byte("")) == 0 {
				insert(v, value, prefix)
				return
			}
			goto ExtensionNode
		}
	ExtensionNode:
		root.right = insertNode(root.right, root.rightLabel, value, gcp, prefix)
		root.rightLabel = []byte("")
		return
	}
}

// InsertBatch makes a new trie
func (trie *MerklePatriciaTrie) InsertBatch(values [][]byte) {
	for _, i := range values {
		trie.Insert(i)
	}
}

// GetRoot gets the root
func (trie *MerklePatriciaTrie) GetRoot() (root *Node) {
	return trie.root

}

func (trie *MerklePatriciaTrie) String() string {
	return trie.root.String()
}
