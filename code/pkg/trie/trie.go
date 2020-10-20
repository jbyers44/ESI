package trie

import (
	"bytes"
	"crypto/sha256"
)

// MerklePatriciaTrie is effectively a wrapper for a pointer to the Root node
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
		trie.root = &Node{[]byte{}, value, leaf, []byte(""), dummy}
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
	return &Node{[]byte{}, labelRemainder, node, valueRemainder, newLeaf}
}

// Private recursive func that does actual insertion logic
func insert(root *Node, value []byte, prefix int) {
	//  Greatest common prefix
	gcp := 0

	// Handle Left pointer
	gcp = greatestCommonPrefix(value[prefix:], root.leftLabel)

	if gcp > 0 {
		switch v := root.left.(type) {
		case *Leaf: // Push leaf down and insert intermediary node with new leaf value on Right
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
			} // The prefix isn't a full match, create a new intermediary node with a new leaf value on Right, old tree on Left
			root.left = insertNode(v, root.leftLabel, value, gcp, prefix)
			root.leftLabel = root.leftLabel[:gcp]
			return
		}
	}

	// Handle Right pointer
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
			} // The prefix isn't a full match, create a new intermediary node with a new leaf value on right, old tree on Left
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

// GenerateHashes generates the merkle Hashes for the tree
func (trie *MerklePatriciaTrie) GenerateHashes() {
	trie.root.hash = Hash(trie.root)
}

// ComputeHash computes the merkle tree Hash for the tree
func Hash(node interface{}) []byte {
	switch n := node.(type) {
	case *Leaf:
		return n.hash

	case *Node:
		h := sha256.New()
		b := append(Hash(n.left), Hash(n.right)...)
		h.Write(b)
		n.hash = h.Sum(nil)
		return n.hash
	}
	println("Hash computation error.\n")
	return []byte{}
}

// Private recursive func that does actual insertion logic
func (trie *MerklePatriciaTrie) InTrie(root *Node, value []byte, prefix int, hashes [][]byte) (bool, [][]byte) {
	//  Greatest common prefix
	gcp := 0

	if len(value) == prefix {
		return true, hashes
	}

	// Handle Left pointer
	gcp = greatestCommonPrefix(value[prefix:], root.leftLabel)
	println(gcp)

	// if gcp > 0 {
	switch v := root.left.(type) {
	case *Leaf: // Push leaf down and insert intermediary node with new leaf value on Right
		if bytes.Compare(v.value, value) == 0 {
			return true, hashes
		}

	case *Node:
		if gcp == len(root.leftLabel) { // If the prefix is a full match, traverse down to the next node
			println(gcp)
			println(string(root.leftLabel))
			switch x := root.right.(type) {
			case *Node:
				hashes = append(hashes, root.hash, x.hash)
			case *Leaf:
				hashes = append(hashes, root.hash, x.hash)
			}
			print("val len: ")
			print(len(value))
			println("prefix+gcp: ")
			print(prefix + gcp)
			// if !(len(value) == prefix+gcp) {
			// 	return trie.InTrie(v, value, prefix+gcp, hashes)
			// }
			return trie.InTrie(v, value, prefix+gcp, hashes)
		}
	}
	// }

	// Handle Right pointer
	gcp = greatestCommonPrefix(value[prefix:], root.rightLabel)

	// if gcp > 0 {

	switch v := root.right.(type) {
	case *Leaf: // Push leaf down and insert intermediary node with new leaf value on Right
		if bytes.Compare(v.value, value) == 0 {
			return true, hashes
		}

	case *Node:
		if gcp == len(root.rightLabel) { // If the prefix is a full match, traverse down to the next node
			switch x := root.left.(type) {
			case *Node:
				hashes = append(hashes, root.hash, x.hash)
			case *Leaf:
				hashes = append(hashes, root.hash, x.hash)
			}
			// if !(len(value) == prefix+gcp) {
			// 	return trie.InTrie(v, value, prefix+gcp, hashes)
			// }
			return trie.InTrie(v, value, prefix+gcp, hashes)
		}
	}
	// }

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
		b := append(Hash(n.left), Hash(n.right)...)
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
	return trie.root

}

func (trie *MerklePatriciaTrie) String() string {
	return trie.root.String()
}
