package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConfirmAction(prompt string) bool {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"

}
