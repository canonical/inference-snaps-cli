package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ErrPromptUnanswerable is returned by PromptYN when the question cannot be
// answered: stdin reached EOF or could not be read. It is not a decline. The
// caller decides what an unanswerable question means, because that depends on
// what has already happened by the time it is asked.
var ErrPromptUnanswerable = errors.New("cannot read a response: stdin is closed or unreadable; use --assume-yes for unattended runs")

// PromptYN prompts the user and returns true for 'y', false for 'n'. It returns
// ErrPromptUnanswerable if no response can be obtained; the boolean is not
// meaningful in that case.
func PromptYN(prompt string, defaultResponse bool) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		if defaultResponse == true {
			fmt.Printf("%s [Y/n] ", prompt) // default is yes
		} else {
			fmt.Printf("%s [y/N] ", prompt) // default is no
		}

		input, err := reader.ReadString('\n')
		eof := errors.Is(err, io.EOF)
		if err != nil && !eof {
			return false, fmt.Errorf("%w: %v", ErrPromptUnanswerable, err)
		}

		// At EOF the read may still carry a final line without a trailing
		// newline. Honour that answer, but never fall back to the default or
		// loop again, since no further input can arrive.
		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y":
			return true, nil
		case "n":
			return false, nil
		case "": // default on empty input
			if !eof {
				return defaultResponse, nil
			}
		default:
			if !eof {
				fmt.Println(`Invalid input. Please enter "y" or "n".`)
				continue
			}
		}

		// Only reachable at EOF: there is no answer and none can arrive.
		fmt.Println()
		return false, ErrPromptUnanswerable
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
	if !assumeYes {
		msg := fmt.Sprintf("Restart %s to apply the changes?", ctx.Snap.InstanceName())
		restart, err := PromptYN(msg, true)
		if err != nil {
			// The configuration has already been written. Skipping the restart
			// silently would leave the stored configuration and the running
			// service diverged while still reporting success, so report it.
			return fmt.Errorf(
				"configuration was changed but %s was not restarted, so it is not yet in effect: %w",
				ctx.Snap.InstanceName(), err,
			)
		}
		if !restart {
			return nil
		}
	}

	if err := ctx.Snap.Restart(); err != nil {
		return fmt.Errorf("restarting snap: %v", err)
	}
	return nil
}
