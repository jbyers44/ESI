package main

import (
	"ESI/pkg/data"
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	println("Output file name:")
	scanner.Scan()
	filename := scanner.Text()

	println("Amount of items:")
	scanner.Scan()
	items := scanner.Text()

	count, err := strconv.Atoi(items)
	if err != nil {
		log.Fatal(err)
	}

	println("Choose between 1. Alphanumeric and 2. ASCII space (includes non-newline whitespace characters and symbols) by entering the number:")
	scanner.Scan()
	choice := scanner.Text()

	if choice == "1" {
		data.GenerateAlphanum(filename, count)
	} else if choice == "2" {
		data.Generate(filename, count)
	} else {
		println("Not a valid choice.")
	}
}
