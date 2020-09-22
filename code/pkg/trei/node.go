package trei

// Node represents a non-terminal node that points to either other nodes or a leaf
type Node struct {
	hash       [32]byte
	leftLabel  []byte
	left       interface{}
	rightLabel []byte
	right      interface{}
}

func (node *Node) SetLeft(left interface{}) {
	node.left = left
}

func (node *Node) SetRight(right interface{}) {
	node.right = right
}
