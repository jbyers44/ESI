package main

import (
	"ESI/pkg/chain"
	"ESI/pkg/helpers"
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

	chain := chain.NewChain()

	var filesData [][]byte

	numFiles := 2
	var filename string
	for i := 0; i < numFiles; i++ {
		content, name := helpers.GetFile()
		filename = name
		filesData = append(filesData, content)
		chain.InsertBatch(filesData)
	}

	trimmed := strings.TrimSuffix(filename, filepath.Ext(filename))
	file, err := os.Create(trimmed + ".block.out.txt")
	check(err)

	defer file.Close()

	for _, block := range chain.GetChain() {
		file.Write([]byte(block.String(false)))
	}

}
