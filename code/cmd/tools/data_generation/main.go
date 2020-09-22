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
	data.Generate(filename, count)
}
