package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfirmAction prompts the user for a message and waits for yes/no confirmation
// it reads user input to stdin and returns true if the user responds with "y" or "yes"
// any other input is treated as a negative response
func ConfirmAction(prompt string) bool {

	// create new buffered reader to read input from terminal
	reader := bufio.NewReader(os.Stdin)

	// print the prompt message with [y/N] suffix
	fmt.Printf("%s [y/N]: ", prompt)

	// read a line of input from the user
	input, _ := reader.ReadString('\n')

	// normalise input, trim whitespace and convert to lower case
	input = strings.ToLower(strings.TrimSpace(input))

	// return true only if input is y or yes
	return input == "y" || input == "yes"
}
