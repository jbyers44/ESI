package trie

import (
	"bytes"
	"crypto/sha256"
)

// MerklePatriciaTrie is effectively a wrapper for a pointer to the Root node
type MerklePatriciaTrie struct {
	Root *Node
}

// NewMerklePatriciaTrie creates a new trie with a root node with two empty string labels
func NewMerklePatriciaTrie() *MerklePatriciaTrie {
	return &MerklePatriciaTrie{nil}
}

// Insert a node into the trie
func (trie *MerklePatriciaTrie) Insert(value []byte) {
	if trie.Root == nil {
		leaf := NewLeaf(value)
		dummy := NewLeaf([]byte(""))
		trie.Root = &Node{[]byte{}, value, leaf, []byte(""), dummy}
	} else {
		insert(trie.Root, value, 0)
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
	return &Node{[]byte{}, labelRemainder, node, valueRemainder, newLeaf}
}

// Private recursive func that does actual insertion logic
func insert(root *Node, value []byte, prefix int) {
	//  Greatest common prefix
	gcp := 0

	// Handle Left pointer
	gcp = greatestCommonPrefix(value[prefix:], root.LeftLabel)

	if gcp > 0 {
		switch v := root.Left.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on Right
			if bytes.Compare(v.value, value) == 0 {
				return
			}
			root.Left = insertNode(v, root.LeftLabel, value, gcp, prefix)
			root.LeftLabel = root.LeftLabel[:gcp]
			return

		case *Node:
			if gcp == len(root.LeftLabel) { // If the prefix is a full match, traverse down to the next node
				insert(v, value, prefix+gcp)
				return
			} // The prefix isn't a full match, create a new intermediary node with a new leaf value on Right, old tree on Left
			root.Left = insertNode(v, root.LeftLabel, value, gcp, prefix)
			root.LeftLabel = root.LeftLabel[:gcp]
			return
		}
	}

	// Handle Right pointer
	gcp = greatestCommonPrefix(value[prefix:], root.RightLabel)

	if gcp > 0 {

		switch v := root.Right.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on Right
			if bytes.Compare(v.value, value) == 0 {
				return
			}
			root.Right = insertNode(v, root.RightLabel, value, gcp, prefix)
			root.RightLabel = root.RightLabel[:gcp]
			return

		case *Node:
			if gcp == len(root.RightLabel) { // If the prefix is a full match, traverse down to the next node
				insert(v, value, prefix+gcp)
				return
			} // The prefix isn't a full match, create a new intermediary node with a new leaf value on Right, old tree on Left
			root.Right = insertNode(v, root.RightLabel, value, gcp, prefix)
			root.RightLabel = root.RightLabel[:gcp]
			return
		}
	} else { // Final case, shares no prefix with either node

		switch v := root.Right.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on Right
			if bytes.Compare(v.value, []byte("")) == 0 {
				root.Right = NewLeaf(value)
				root.RightLabel = value
				return
			}
			goto ExtensionNode

		case *Node:
			if bytes.Compare(root.RightLabel, []byte("")) == 0 {
				insert(v, value, prefix)
				return
			}
			goto ExtensionNode
		}
	ExtensionNode:
		root.Right = insertNode(root.Right, root.RightLabel, value, gcp, prefix)
		root.RightLabel = []byte("")
		return
	}
}

// GenerateHashes generates the merkle Hashes for the tree
func (trie *MerklePatriciaTrie) GenerateHashes() {
	trie.Root.Hash = Hash(trie.Root)
}

// ComputeHash computes the merkle tree Hash for the tree
func Hash(node interface{}) []byte {
	switch n := node.(type) {
	case *Leaf:
		return n.Hash

	case *Node:
		h := sha256.New()
		b := append(Hash(n.Left), Hash(n.Right)...)
		h.Write(b)
		n.Hash = h.Sum(nil)
		return n.Hash
	}
	println("Hash computation error.\n")
	return []byte{}
}

// Private recursive func that does actual insertion logic
func (trie *MerklePatriciaTrie) InTrie(root *Node, value []byte, prefix int, hashes [][]byte) (bool, [][]byte) {
	//  Greatest common prefix
	gcp := 0

	// Handle Left pointer
	gcp = greatestCommonPrefix(value[prefix:], root.LeftLabel)

	if gcp > 0 {
		switch v := root.Left.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on Right
			if bytes.Compare(v.value, value) == 0 {
				return true, hashes
			}

		case *Node:
			if gcp == len(root.LeftLabel) { // If the prefix is a full match, traverse down to the next node
				switch x := root.Right.(type) {
				case *Node:
					hashes = append(hashes, root.Hash, x.Hash)
					return trie.InTrie(v, value, prefix+gcp, hashes)
				case *Leaf:
					hashes = append(hashes, root.Hash, x.Hash)
				}
			}
		}
	}

	// Handle Right pointer
	gcp = greatestCommonPrefix(value[prefix:], root.RightLabel)

	if gcp > 0 {

		switch v := root.Right.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on Right
			if bytes.Compare(v.value, value) == 0 {
				return true, hashes
			}

		case *Node:
			if gcp == len(root.RightLabel) { // If the prefix is a full match, traverse down to the next node
				switch x := root.Left.(type) {
				case *Node:
					hashes = append(hashes, root.Hash, x.Hash)
					return trie.InTrie(v, value, prefix+gcp, hashes)
				case *Leaf:
					hashes = append(hashes, root.Hash, x.Hash)
				}
			}
		}
	}

	return false, nil
}

// InsertBatch makes a new trie
func (trie *MerklePatriciaTrie) InsertBatch(values [][]byte) {
	for _, i := range values {
		trie.Insert(i)
	}
}

func (trie *MerklePatriciaTrie) Validate() bool {
	if trie.root == nil {
		return true
	} else {
		return validate(trie.root)
	}
}

func validate(node interface{}) bool {

	switch n := node.(type) {
	case *Leaf:
		h := sha256.New()
		h.Write(n.value)
		if bytes.Compare(n.hash, h.Sum(nil)) != 0 {
			return false
		}

	case *Node:
		h := sha256.New()
		b := append(hash(n.left), hash(n.right)...)
		h.Write(b)
		if bytes.Compare(n.hash, h.Sum(nil)) != 0 {
			return false
		}
		if n.left != nil {
			return validate(n.left)
		}
		if n.right != nil {
			return validate(n.right)
		}
	}

	return true
}

// GetRoot gets the root
func (trie *MerklePatriciaTrie) GetRoot() (Root *Node) {
	return trie.Root

}

func (trie *MerklePatriciaTrie) String() string {
	return trie.Root.String()
}
