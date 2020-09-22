package main

import (
	"ESI/pkg/trie"
	"bufio"
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

	scanner := bufio.NewScanner(os.Stdin)
	println("Please enter your input filename (which should be placed in the same directory as hw3.exe)")

	scanner.Scan()
	input := scanner.Text()

	content, err := ioutil.ReadFile(input)
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
