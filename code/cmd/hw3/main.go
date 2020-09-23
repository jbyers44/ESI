package main

import (
	"ESI/pkg/trie"
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	content, filename := getInput()

	// Convert []byte to string and print to screen
	byteStrings := bytes.Split(content, []byte("\n"))

	sort.Slice(byteStrings, func(i, j int) bool {
		return bytes.Compare(byteStrings[i], byteStrings[j]) < 0
	})

	mpt := trie.NewMerklePatriciaTrie()
	mpt.InsertBatch(byteStrings)
	mpt.GenerateHashes()

	trimmed := strings.TrimSuffix(filename, filepath.Ext(filename))
	file, err := os.Create(trimmed + ".out.txt")
	check(err)

	defer file.Close()
	file.Write([]byte(mpt.String()))
}

func getInput() ([]byte, string) {
	scanner := bufio.NewScanner(os.Stdin)
	var input string

	println("Please enter your input filename (which should be placed in the same directory as hw3.exe)")

	for {
		scanner.Scan()
		input = scanner.Text()

		content, err := ioutil.ReadFile(input)
		if err != nil {
			println("Invalid filename, please input a valid file.")
		} else {
			return content, input
		}
	}
}
