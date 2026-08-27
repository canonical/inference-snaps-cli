package common

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
)

func withStdin(input string, fn func()) {
	originalStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	if _, err := w.WriteString(input); err != nil {
		panic(err)
	}
	_ = w.Close()

	os.Stdin = r
	defer func() {
		os.Stdin = originalStdin
		_ = r.Close()
	}()

	fn()
}

func printToStdout(a any) {
	fmt.Printf("-> %v\n", a)
}

func mustPromptYN(prompt string, defaultResponse bool) bool {
	answer, err := PromptYN(prompt, defaultResponse)
	if err != nil {
		panic(err)
	}

	return answer
}

func ExamplePromptYN_defaultYes() {
	withStdin("\n", func() {
		printToStdout(mustPromptYN("Proceed?", true))
	})

	// Output:
	// Proceed? [Y/n] -> true
}

func ExamplePromptYN_invalidThenNo() {
	withStdin("maybe\nn\n", func() {
		printToStdout(mustPromptYN("Proceed?", true))
	})

	// Output:
	// Proceed? [Y/n] Invalid input. Please enter "y" or "n".
	// Proceed? [Y/n] -> false
}

func TestPromptYNUnanswerable(t *testing.T) {
	// An unanswerable prompt is not a decline. Callers need to tell the two
	// apart, because for a prompt asked after state has already changed,
	// silently declining leaves that state diverged while reporting success.
	tests := map[string]struct {
		stdin           string
		defaultResponse bool
	}{
		"closed stdin, default yes": {stdin: "", defaultResponse: true},
		"closed stdin, default no":  {stdin: "", defaultResponse: false},
		"invalid answer then EOF":   {stdin: "maybe\n", defaultResponse: true},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			withStdin(testData.stdin, func() {
				_, err := PromptYN("Proceed?", testData.defaultResponse)
				if !errors.Is(err, ErrPromptUnanswerable) {
					t.Errorf("PromptYN returned error %v, expected ErrPromptUnanswerable", err)
				}
			})
		})
	}
}

func ExamplePromptYN_finalLineWithoutNewline() {
	withStdin("y", func() {
		printToStdout(mustPromptYN("Proceed?", false))
	})

	// Output:
	// Proceed? [y/N] -> true
}

func ExamplePromptlnEnter() {
	withStdin("\n", func() {
		printToStdout(PromptlnEnter("continue"))
	})

	// Output:
	// Press [Enter] to continue, or [Ctrl+C] to abort. -> true
}

func TestPromptRestartToApplyChanges(t *testing.T) {
	// The configuration is already written by the time this prompt is asked, so
	// an unanswerable prompt must not be reported as success: that would leave
	// the stored configuration and the running service diverged with exit 0.
	t.Run("unanswerable prompt is an error", func(t *testing.T) {
		ctx := &Context{Snap: snap.Mock()}

		var err error
		withStdin("", func() {
			err = PromptRestartToApplyChanges(ctx, false)
		})

		if err == nil {
			t.Fatal("expected an error when the restart prompt cannot be answered, got none")
		}
		if !errors.Is(err, ErrPromptUnanswerable) {
			t.Errorf("error %v does not wrap ErrPromptUnanswerable", err)
		}
		if !strings.Contains(err.Error(), "was not restarted") {
			t.Errorf("error %q does not say the service was not restarted", err)
		}
	})

	t.Run("declining is not an error", func(t *testing.T) {
		ctx := &Context{Snap: snap.Mock()}

		var err error
		withStdin("n\n", func() {
			err = PromptRestartToApplyChanges(ctx, false)
		})

		if err != nil {
			t.Errorf("declining the restart returned %v, expected no error", err)
		}
	})

	t.Run("assume-yes restarts without reading stdin", func(t *testing.T) {
		ctx := &Context{Snap: snap.Mock()}

		var err error
		withStdin("", func() {
			err = PromptRestartToApplyChanges(ctx, true)
		})

		if err != nil {
			t.Errorf("--assume-yes returned %v, expected no error", err)
		}
	})
}
