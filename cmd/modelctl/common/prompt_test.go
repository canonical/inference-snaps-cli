package common

import (
	"fmt"
	"os"
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

func ExamplePromptYN_defaultYes() {
	withStdin("\n", func() {
		printToStdout(PromptYN("Proceed?", true))
	})

	// Output:
	// Proceed? [Y/n] -> true
}

func ExamplePromptYN_invalidThenNo() {
	withStdin("maybe\nn\n", func() {
		printToStdout(PromptYN("Proceed?", true))
	})

	// Output:
	// Proceed? [Y/n] Invalid input. Please enter "y" or "n".
	// Proceed? [Y/n] -> false
}

func ExamplePromptYN_eofDeclines() {
	withStdin("", func() {
		printToStdout(PromptYN("Proceed?", true))
	})

	// Output:
	// Proceed? [Y/n] No input available to answer the prompt. Assuming "n"; use --assume-yes to bypass confirmation prompts.
	// -> false
}

func ExamplePromptYN_invalidThenEOF() {
	withStdin("maybe\n", func() {
		printToStdout(PromptYN("Proceed?", true))
	})

	// Output:
	// Proceed? [Y/n] Invalid input. Please enter "y" or "n".
	// Proceed? [Y/n] No input available to answer the prompt. Assuming "n"; use --assume-yes to bypass confirmation prompts.
	// -> false
}

func ExamplePromptYN_finalLineWithoutNewline() {
	withStdin("y", func() {
		printToStdout(PromptYN("Proceed?", false))
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
