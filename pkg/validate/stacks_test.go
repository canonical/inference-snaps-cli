package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStackFiles(t *testing.T) {
	stacksDir := "../../test_data/stacks"

	entries, err := os.ReadDir(stacksDir)
	if err != nil {
		t.Fatalf("Failed reading stacks dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			stack := entry.Name()
			stackPath := filepath.Join(stacksDir, stack, "stack.yaml")
			t.Run(stack, func(t *testing.T) {
				err = Stack(stackPath)
				if err != nil {
					t.Fatalf("%s: %v", stack, err)
				}
			})
		}
	}
}

func TestStackEmpty(t *testing.T) {
	data := ""
	err := validateStackYaml([]byte(data))
	if err == nil {
		t.Fatal("Empty yaml should fail")
	}
}

func TestStackInvalidField(t *testing.T) {
	data := "invalid-field: test"
	err := validateStackYaml([]byte(data))
	if err == nil {
		t.Fatal("Unknown field should fail")
	}
}
