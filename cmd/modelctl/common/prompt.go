package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// PromptYN prompts the user and returns true for 'y', false for 'n'.
func PromptYN(prompt string, defaultResponse bool) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		if defaultResponse == true {
			fmt.Printf("%s [Y/n] ", prompt) // default is yes
		} else {
			fmt.Printf("%s [y/N] ", prompt) // default is no
		}

		input, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Printf("Error reading input: %v\n", err)
			return false
		}
		// On EOF, input may still hold a final line without a trailing
		// newline; honor it, but never fall back to the default or
		// retry, since no further input can arrive.
		eof := err != nil

		input = strings.ToLower(strings.TrimSpace(input))
		switch input {
		case "": // default on empty input
			if !eof {
				return defaultResponse
			}
		case "Y", "y":
			return true
		case "N", "n":
			return false
		default:
			if !eof {
				fmt.Println(`Invalid input. Please enter "y" or "n".`)
			}
		}

		if eof {
			fmt.Println(`No input available to answer the prompt. Assuming "n"; use --assume-yes to bypass confirmation prompts.`)
			return false
		}
	}
}

// PromptlnEnter prompts the user for Enter in a new line
func PromptlnEnter(action string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Press [Enter] to %s, or [Ctrl+C] to abort. ", action)

	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nError reading input: %v\n", err)
		return false
	}

	return true
}

func PromptRestartToApplyChanges(ctx *Context, assumeYes bool) error {
	msg := fmt.Sprintf("Restart %s to apply the changes?", ctx.Snap.InstanceName())
	if assumeYes || PromptYN(msg, true) {
		if err := ctx.Snap.Restart(); err != nil {
			return fmt.Errorf("restarting snap: %v", err)
		}
	}
	return nil
}
