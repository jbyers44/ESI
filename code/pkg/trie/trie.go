package trie

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"time"
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

// Hash computes the merkle tree Hash for the tree
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

// InTrie checks if value is in trie, returns proof if so
func (trie *MerklePatriciaTrie) InTrie(value []byte) (bool, [][]byte) {
	if trie.root == nil {
		return false, nil
	}
	var blank [][]byte
	result, trieHashes := inTrie(trie.root, value, 0, blank)

	// Flip the order
	for j := 0; j < len(trieHashes)/2; j++ {
		k := len(trieHashes) - j - 1
		trieHashes[j], trieHashes[k] = trieHashes[k], trieHashes[j]
	}

	// Append the root
	trieHashes = append(trieHashes, trie.root.hash)
	return result, trieHashes
}

// InTrie checks if value is in trie, returns proof if so
func inTrie(root *Node, value []byte, prefix int, hashes [][]byte) (bool, [][]byte) {
	//  Greatest common prefix
	gcp := 0

	// Handle left label
	gcp = greatestCommonPrefix(value[prefix:], root.leftLabel)

	if gcp == len(root.leftLabel) {
		switch v := root.left.(type) {
		// If root.left is our value, append its neighboring hash, then its hash, and return
		case *Leaf:
			if bytes.Compare(v.value, value) == 0 {
				switch n := root.right.(type) {
				case *Node:
					hashes = append(hashes, n.hash, v.hash)
				case *Leaf:
					hashes = append(hashes, n.hash, v.hash)
				}
				return true, hashes
			} else if bytes.Compare(root.leftLabel, []byte("")) == 0 {
				break
			}
			return false, nil

		// If root.left is a node, append its neighboring hash, its hash, and traverse to it
		case *Node:
			switch n := root.right.(type) {
			case *Node:
				hashes = append(hashes, n.hash, v.hash)
			case *Leaf:
				hashes = append(hashes, n.hash, v.hash)
			}
			return inTrie(v, value, prefix+gcp, hashes)
		}
	}

	// Handle right label
	gcp = greatestCommonPrefix(value[prefix:], root.rightLabel)

	if gcp == len(root.rightLabel) {
		switch v := root.right.(type) {
		// If root.right is our value, append its neighboring hash, then its hash, and return
		case *Leaf:
			if bytes.Compare(v.value, value) == 0 {
				switch n := root.left.(type) {
				case *Node:
					hashes = append(hashes, n.hash, v.hash)
				case *Leaf:
					hashes = append(hashes, n.hash, v.hash)
				}
				return true, hashes
			}
			return false, nil

		// If root.right is a node, append its neighboring hash, its hash, and traverse to it
		case *Node:
			switch n := root.left.(type) {
			case *Node:
				hashes = append(hashes, n.hash, v.hash)
			case *Leaf:
				hashes = append(hashes, n.hash, v.hash)
			}
			return inTrie(v, value, prefix+gcp, hashes)
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

// Validate validates the trie for hash completeness
func (trie *MerklePatriciaTrie) Validate() bool {
	if trie.root == nil {
		return true
	}
	return validate(trie.root)
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

// Corrupt randomly chooses a block in the trie and corrupts its hash (sets it to the hash of the existing hash)
func (trie *MerklePatriciaTrie) Corrupt() {
	if trie.root == nil {
		return
	}
	corrupt(trie.root)
}

func corrupt(node interface{}) {
	rand.Seed(time.Now().UnixNano())

	switch n := node.(type) {
	case *Leaf:
		h := sha256.New()
		h.Write(n.hash)
		n.hash = h.Sum(nil)

	case *Node:
		flip := rand.Intn(2)
		switch flip {
		case 0:
			corrupt(n.left)
		case 1:
			corrupt(n.right)
		}
	}
}

// GetRoot gets the root
func (trie *MerklePatriciaTrie) GetRoot() (Root *Node) {
	return trie.root

}

func (trie *MerklePatriciaTrie) String() string {
	return trie.root.String()
}
