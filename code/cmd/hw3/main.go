package main

import (
	"ESI/pkg/helpers"
	"ESI/pkg/trie"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	content, filename := helpers.GetFile()

	mpt := trie.NewMerklePatriciaTrie()
	mpt.InsertBatch(helpers.SplitBytes(content))
	mpt.GenerateHashes()

	trimmed := strings.TrimSuffix(filename, filepath.Ext(filename))
	file, err := os.Create(trimmed + ".out.txt")
	check(err)

	defer file.Close()
	file.Write([]byte(mpt.String()))

	x, _ := mpt.InTrie(mpt.GetRoot(), []byte("banana"), 0, [][]byte{})
	println(x)
}
