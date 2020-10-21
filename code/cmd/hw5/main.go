package main

import (
	"ESI/pkg/chain"
	"ESI/pkg/data"
	"ESI/pkg/helpers"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)
var currentChain *chain.Chain
var currentChainName string = ""

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mainmenu()
}

func mainmenu() {
	fmt.Fprintf(os.Stdout, "\nChoose one of the following options by entering the corresponding number:\n")
	fmt.Fprintf(os.Stdout, "	1. Create a new chain from 1 or more files (lists of strings).\n")
	fmt.Fprintf(os.Stdout, "	2. Randomly generate a new chain.\n")
	fmt.Fprintf(os.Stdout, "	3. Quit.\n")

	scanner.Scan()
	input := scanner.Text()

	switch input {
	case "1":
		loadchain()
	case "2":
		genchain()
	case "3":
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stdout, "That's not a valid option.\n")
		mainmenu()
	}
}

func loadchain() {
	currentChain = chain.NewChain()

	fmt.Fprintf(os.Stdout, "\nEnter the number of files you want to input.\n")
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
		fmt.Fprintf(os.Stdout, "File #%d loaded successfully.\n", i+1)
	}

	currentChain.InsertBatch(contents)
	currentChainName = strings.TrimSuffix(filename, filepath.Ext(filename))

	fmt.Fprintf(os.Stdout, "\nChain successfully loaded, bringing you to the chain management menu.\n")
	chainmenu()
}

func genchain() {
	currentChain = chain.NewChain()

	fmt.Fprintf(os.Stdout, "\nEnter the number of blocks you want in the chain.\n")
	scanner.Scan()
	numBlocks, err := strconv.Atoi(scanner.Text())
	check(err)

	fmt.Fprintf(os.Stdout, "\nEnter the amount of values per block.\n")
	scanner.Scan()
	numValues, err := strconv.Atoi(scanner.Text())
	check(err)

	contents := make(map[string][]byte)
	currentChainName = "random"
	for i := 0; i < numBlocks; i++ {
		filename := "random" + strconv.Itoa(i)
		data.GenerateAlphanum(filename+".tmp", numValues)
		content, err := ioutil.ReadFile(filename + ".tmp")
		check(err)
		contents[filename] = content
		err = os.Remove(filename + ".tmp")
		check(err)
	}

	currentChain.InsertBatch(contents)

	fmt.Fprintf(os.Stdout, "\nChain successfully generated, bringing you to the chain management menu.\n")
	chainmenu()
}

func chainmenu() {
	fmt.Fprintf(os.Stdout, "\nChoose one of the following options by entering the corresponding number:\n")
	fmt.Fprintf(os.Stdout, "	1. Output the chain to a file.\n")
	fmt.Fprintf(os.Stdout, "	2. Corrupt a block (change a random hash in a random trie).\n")
	fmt.Fprintf(os.Stdout, "	3. Search for a value, returning a proof of membership if found.\n")
	fmt.Fprintf(os.Stdout, "	4. Validate the chain.\n")
	fmt.Fprintf(os.Stdout, "	5. Make a new chain (main menu).\n")
	fmt.Fprintf(os.Stdout, "	6. Quit.\n")

	scanner.Scan()
	input := scanner.Text()

	switch input {
	case "1":
		outputchain()
	case "2":
		corruptchain()
	case "3":
		searchchain()
	case "4":
		validatechain()
	case "5":
		mainmenu()
	case "6":
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stdout, "That's not a valid option.\n")
		chainmenu()
	}
}

func outputchain() {
	fmt.Fprintf(os.Stdout, "\nChoose one of the following options by entering the corresponding number:\n")
	fmt.Fprintf(os.Stdout, "	1. Output with MPT printed.\n")
	fmt.Fprintf(os.Stdout, "	2. Output with no MPT.\n")
	fmt.Fprintf(os.Stdout, "	3. Cancel.\n")

	scanner.Scan()
	input := scanner.Text()

	file, err := os.Create(currentChainName + ".block.out")
	check(err)
	defer file.Close()

	switch input {
	case "1":
		file.Write([]byte(currentChain.String(true)))
		fmt.Fprintf(os.Stdout, "\nSuccessfully output chain.\n")
		chainmenu()
	case "2":
		file.Write([]byte(currentChain.String(false)))
		fmt.Fprintf(os.Stdout, "\nSuccessfully output chain.\n")
		chainmenu()
	case "3":
		chainmenu()
	default:
		fmt.Fprintf(os.Stdout, "\nThat's not a valid option.\n")
		chainmenu()
	}
}

func corruptchain() {
	currentChain.Corrupt()
	fmt.Fprintf(os.Stdout, "\nA random node in a random block has been corrupted.\n")
	chainmenu()
}

func searchchain() {
	fmt.Fprintf(os.Stdout, "\nEnter the value you want to search:\n")
	scanner.Scan()
	input := scanner.Text()

	result, trieHashes, blockHashes := currentChain.InChain(input)

	if result {
		fmt.Fprintf(os.Stdout, "The value was successfully found in the chain.\n")
		fmt.Fprintf(os.Stdout, "Merkle proof from value to merkle root: \n[\n")
		for _, i := range trieHashes {
			fmt.Fprintf(os.Stdout, "%x,\n", i)
		}
		fmt.Fprintf(os.Stdout, "]\n\nBlock hashes from value's block to most recent: \n[\n")
		for _, i := range blockHashes {
			fmt.Fprintf(os.Stdout, "%x,\n", i)
		}
		fmt.Fprintf(os.Stdout, "]\n\n")
		fmt.Fprintf(os.Stdout, "Press 'Enter' to return to the chain management menu...")
		fmt.Scanln()
		chainmenu()
	} else {
		fmt.Fprintf(os.Stdout, "The value was not found in the chain.\n")
		fmt.Fprintf(os.Stdout, "Press 'Enter' to return to the chain management menu...")
		fmt.Scanln()
		chainmenu()
	}
}

func validatechain() {
	fmt.Fprintf(os.Stdout, "\nValidating chain...\n")
	result := currentChain.Validate()

	if result {
		fmt.Fprintf(os.Stdout, "The chain is currently valid.\n")
	} else {
		fmt.Fprintf(os.Stdout, "The chain is currently invalid, there is a corrupt block.\n")
	}
	chainmenu()
}
