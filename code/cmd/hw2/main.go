package main

import (
	"ESI/pkg/chain"
	"ESI/pkg/helpers"
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	chain := chain.NewChain()

	scanner := bufio.NewScanner(os.Stdin)

	println("Please enter how many files you want to input.")

	scanner.Scan()
	numFiles, _ := strconv.Atoi(scanner.Text())

	var filesData [][]byte
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