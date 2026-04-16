package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompt prompts the user and returns true for 'y', false for 'n'.
func PromptYN(prompt string, defaultResponse bool) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		if defaultResponse == true {
			fmt.Printf("%s [Y/n] ", prompt) // default is yes
		} else {
			fmt.Printf("%s [y/N] ", prompt) // default is no
		}

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.ToLower(strings.TrimSpace(input))
		switch input {
		case "": // default on empty input
			return defaultResponse
		case "Y", "y":
			return true
		case "N", "n":
			return false
		default:
			fmt.Println(`Invalid input. Please enter "y" or "n".`)
		}
	}
}

// PromptEnterln prompts the user for Enter in a new line
func PromptlnEnter(action string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Press [Enter] to %s, or [q] to abort. ", action)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nError reading input: %v\n", err)
		return false
	}

	if strings.TrimSpace(input) == "q" {
		return false
	}

	return true
}
