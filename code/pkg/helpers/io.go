package helpers

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"sort"
)

// GetFile prompts the user for a valid filename until they enter one correctly, returns the content of the file along with the input that was successful (the filename)
func GetFile() ([]byte, string) {
	scanner := bufio.NewScanner(os.Stdin)
	var input string

	println("Enter an input filepath:")

	for {
		scanner.Scan()
		input = scanner.Text()

		content, err := ioutil.ReadFile(input)
		if err != nil {
			println("Invalid filepath, please input a valid file.")
		} else {
			return content, input
		}
	}
}

// SplitBytes splits a byteslice on newline characters and sorts them alphabetically
func SplitBytes(content []byte) [][]byte {
	byteStrings := bytes.Split(content, []byte("\n"))

	sort.Slice(byteStrings, func(i, j int) bool {
		return bytes.Compare(byteStrings[i], byteStrings[j]) < 0
	})

	return byteStrings
}
