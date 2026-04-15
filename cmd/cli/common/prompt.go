package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfirmationPrompt prompts the user and returns true for 'y', false for 'n'.
// It defaults to true if the user presses Enter without a response.
func ConfirmationPrompt(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [Y/n] ", prompt) // default to yes

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.ToLower(strings.TrimSpace(input))
		switch input {
		case "": // default to yes on empty input
			return true
		case "Y", "y":
			return true
		case "N", "n":
			return false
		default:
			fmt.Println(`Invalid input. Please enter "y" or "n".`)
		}
	}
}

// ConfirmationPromptEnter prompts the user for Enter
func ConfirmationPromptEnter(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s Press [Enter] to continue... ", prompt)

	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nError reading input: %v\n", err)
		return false
	}
	return true
}
