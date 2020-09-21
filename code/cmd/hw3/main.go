package main

import (
	"ESI/pkg/trei"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

func main() {

	content, err := ioutil.ReadFile("code/cmd/hw3/testInput.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	byteStrings := bytes.Split(content, []byte("\n"))

	sort.Slice(byteStrings, func(i, j int) bool {
		return bytes.Compare(byteStrings[i], byteStrings[j]) < 0
	})

	// Simply to test reading works
	for i, s := range byteStrings {
		fmt.Println(i, string(s))
	}

	merkleTree := new(trei.MerklePatriciaTrie)
	merkleTree.NewTrie(byteStrings)

}
