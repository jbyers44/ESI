package main

import(
	"os"
	"fmt"
	"strconv"
)

//write data from all nodes to text file
func printTree(node *Node, file *os.File) {
	data := [][][]byte{}
	dataPtr := &data

	getDataFromAllNodes(node, dataPtr, 1)

	sortData(dataPtr)

	for i := 0; i < len(data); i++ {
		if len(data[i]) > 3 {
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
		} else {
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

//store fields from leaf nodes in 
func getLeafData(leaf *Leaf, data *[][][]byte, printID int) {
	leafData := [][]byte{}
	str := strconv.Itoa(printID)
	leafData = append(leafData, []byte(str))
	leafData = append(leafData, leaf.value)
	leafData = append(leafData, leaf.hash[:])

	*data = append(*data, leafData)
}

//sorts array of node data
func sortData(data *[][][]byte) {
	fmt.Println(*data)
	minIndex := 0
	min := 0
	curr := 0

	for i := 0; i < len(*data) - 1; i++ {
		minIndex = i
		for j := i + 1; j < len(*data); j++ {
			curr, _ = strconv.Atoi(string((*data)[j][0]))
			min, _ = strconv.Atoi(string((*data)[minIndex][0]))

			if curr < min {
				minIndex = j
			}
		}

		(*data)[i], (*data)[minIndex] = (*data)[minIndex], (*data)[i]
	}
}