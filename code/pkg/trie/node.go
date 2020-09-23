package trie

import (
	"bytes"
	"fmt"
	"sort"
)

// Node represents a non-terminal node that points to either other nodes or a leaf
type Node struct {
	hash       []byte
	leftLabel  []byte
	left       interface{}
	rightLabel []byte
	right      interface{}
}

func (node *Node) String() string {
	m := make(map[int]string)
	stringify(node, m)

	nodeIds := make([]int, len(m))
	i := 0
	for id := range m {
		nodeIds[i] = id
		i++
	}
	sort.Ints(nodeIds)

	var b bytes.Buffer
	for _, id := range nodeIds {
		fmt.Fprintf(&b, "%s\n\n", m[id])
	}
	return b.String()
}

func stringify(node interface{}, nodeMap map[int]string) {
	switch n := node.(type) {
	case *Leaf:
		var b bytes.Buffer
		fmt.Fprintf(&b, "[printID]		%d\n", len(nodeMap)+1)
		fmt.Fprintf(&b, "[value]			'%s'\n", n.value)
		fmt.Fprintf(&b, "[hash]			%x\n", n.hash)
		nodeMap[len(nodeMap)+1] = b.String()
		return

	case *Node:
		var b bytes.Buffer
		fmt.Fprintf(&b, "[printId]		%d\n", len(nodeMap)+1)
		fmt.Fprintf(&b, "[left-child]	%d\n", len(nodeMap)+2)
		fmt.Fprintf(&b, "[left-label]	'%s'\n", n.leftLabel)
		fmt.Fprintf(&b, "[hash]			%x\n", n.hash)
		fmt.Fprintf(&b, "[right-label]	'%s'\n", n.rightLabel)
		fmt.Fprintf(&b, "[right-child]	%d\n", len(nodeMap)+3)
		nodeMap[len(nodeMap)+1] = b.String()

		stringify(n.left, nodeMap)
		stringify(n.right, nodeMap)
		return
	}
}
