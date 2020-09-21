package trei

// Node represents a non-terminal node that points to either other nodes or a leaf
type Node struct {
	hash       [32]byte
	leftLabel  []byte
	left       interface{}
	rightLabel []byte
	right      interface{}
}
