package main

import (
	"ESI/pkg/trie"
	"bufio"
	"bytes"
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

	println("Enter the name of the file in which to save the printed tree (the program will automatically append the appropriate file extension)")
	scanner.Scan()
	input = scanner.Text()

	// Convert []byte to string and print to screen
	byteStrings := bytes.Split(content, []byte("\n"))

	sort.Slice(byteStrings, func(i, j int) bool {
		return bytes.Compare(byteStrings[i], byteStrings[j]) < 0
	})

	mpt := trie.NewMerklePatriciaTrie()
	mpt.InsertBatch(byteStrings)

	file, err := os.Create(input + ".out.txt")
	check(err)

	defer file.Close()
	file.Write([]byte(mpt.String()))
}
