package main

import (
	"ESI/pkg/chain"
	"ESI/pkg/helpers"
	"bufio"
	"fmt"
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

	println("Please enter the number of files you want to input.")

	scanner.Scan()
	numFiles, err := strconv.Atoi(scanner.Text())
	check(err)

	contents := make(map[string][]byte)
	var filename string
	for i := 0; i < numFiles; i++ {
		content, name := helpers.GetFile()
		contents[name] = content
		if i == 0 {
			filename = name
		}
		fmt.Fprintf(os.Stdout, "File #%d loaded successfully.\n\n", i+1)
	}

	chain.InsertBatch(contents)
	trimmed := strings.TrimSuffix(filename, filepath.Ext(filename))
	file, err := os.Create(trimmed + ".block.out.txt")
	check(err)

	defer file.Close()

	file.Write([]byte(chain.String(true)))

	println(chain.Validate())
}
