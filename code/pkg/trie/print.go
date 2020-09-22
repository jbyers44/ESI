package trie

import (
	"fmt"
	"os"
	"strconv"
)

//check error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// WriteTree data from all nodes to text file
func WriteTree(node *Node, file *os.File) {
	data := [][][]byte{}
	dataPtr := &data

	getDataFromAllNodes(node, dataPtr, 1)
	fmt.Println(data)

	for i := 0; i < len(data); i++ {
		if len(data[i]) > 3 { //write internal nodes
			file.Write([]byte("printID: "))
			file.Write(data[i][0])
			file.Write([]byte("\n"))
			file.Write([]byte("hash: "))
			file.Write(data[i][1])
			file.Write([]byte("\n"))
			file.Write([]byte("left edge: "))
			file.Write(data[i][2])
			file.Write([]byte("\n"))
			file.Write([]byte("left child printID: "))
			file.Write(data[i][3])
			file.Write([]byte("\n"))
			file.Write([]byte("right edge: "))
			file.Write(data[i][4])
			file.Write([]byte("\n"))
			file.Write([]byte("right child printID: "))
			file.Write(data[i][5])
			file.Write([]byte("\n\n\n"))
		} else { //write leaf nodes
			file.Write([]byte("printID: "))
			file.Write(data[i][0])
			file.Write([]byte("\n"))
			file.Write([]byte("value: "))
			file.Write(data[i][1])
			file.Write([]byte("\n"))
			file.Write([]byte("hash: "))
			file.Write(data[i][2])
			file.Write([]byte("\n\n\n"))
		}
	}
}

//recursively store all fields from all nodes in data array
func getDataFromAllNodes(node *Node, data *[][][]byte, printID int) {
	_, lLeaf := node.left.(Leaf)
	_, rLeaf := node.right.(Leaf)

	nodeData := [][]byte{}
	str := strconv.Itoa(printID)
	nodeData = append(nodeData, []byte(str))
	nodeData = append(nodeData, node.hash[:])
	nodeData = append(nodeData, []byte(node.leftLabel))
	str = strconv.Itoa(2 * printID)
	nodeData = append(nodeData, []byte(str))
	nodeData = append(nodeData, []byte(node.rightLabel))
	str = strconv.Itoa(2*printID + 1)
	nodeData = append(nodeData, []byte(str))

	*data = append(*data, nodeData)

	if lLeaf {
		lLeaf := node.left.(Leaf)
		getLeafData(&lLeaf, data, 2*printID)
	} else {
		lNode := node.left.(Node)
		getDataFromAllNodes(&lNode, data, 2*printID)
	}

	if rLeaf {
		rLeaf := node.right.(Leaf)
		getLeafData(&rLeaf, data, 2*printID+1)
	} else {
		rNode := node.right.(Node)
		getDataFromAllNodes(&rNode, data, 2*printID+1)
	}
}

//store fields from leaf nodes in
func getLeafData(leaf *Leaf, data *[][][]byte, printID int) {
	leafData := [][]byte{}
	str := strconv.Itoa(printID)
	leafData = append(leafData, []byte(str))
	leafData = append(leafData, leaf.value)
	leafData = append(leafData, leaf.hash[:])

	*data = append(*data, leafData)
}
