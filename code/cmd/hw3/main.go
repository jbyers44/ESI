package main

import (
	"ESI/pkg/trie"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	content, err := ioutil.ReadFile("./testInput.txt")
	check(err)

	// Convert []byte to string and print to screen
	byteStrings := bytes.Split(content, []byte("\n"))

	sort.Slice(byteStrings, func(i, j int) bool {
		return bytes.Compare(byteStrings[i], byteStrings[j]) < 0
	})

	// Simply to test reading works
	for i, s := range byteStrings {
		fmt.Println(i, string(s))
	}

	merkleTree := trie.NewMerklePatriciaTrie()
	merkleTree.NewTrie(byteStrings)

	file, err := os.Create("output.txt")
	check(err)

	defer file.Close()
	trie.PrintTrie(merkleTree.GetRoot(), file)
}
