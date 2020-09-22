package data

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Generate creates a file in the script directory with name filename consisting of {count} amount of random, newline separated strings
func Generate(filename string, count int) {
	file, err := os.Create(filename)
	check(err)

	defer file.Close()

	writer := bufio.NewWriter(file)

	rand.Seed(time.Now().UnixNano())

	buffer := make([]byte, 100)
	for i := 0; i < count; i++ {
		rand.Read(buffer)

		for _, byte := range buffer {
			shifted := (byte >> 1)
			fmt.Fprintf(writer, "%c", rune(shifted))
		}
		fmt.Fprintf(writer, "\n")
	}
}
