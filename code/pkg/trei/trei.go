package trei

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
		insert(trie.root, value)
	}
}

// Private recursive func that does actual insertion logic
func insert(root *Node, value []byte) {

}
