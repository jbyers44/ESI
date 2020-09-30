package helpers

import (
	"bufio"
	"io/ioutil"
	"os"
)

// GetFile prompts the user for a valid filename until they enter one correctly, returns the content of the file along with the input that was successful (the filename)
func GetFile() ([]byte, string) {
	scanner := bufio.NewScanner(os.Stdin)
	var input string

	println("Please enter your input filepath.")

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
