package main

import (
	"ESI/pkg/helpers"
	"ESI/pkg/trie"
	"bytes"
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
	content, filename := helpers.GetFile()

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
