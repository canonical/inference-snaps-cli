package commands

import (
	"os"
	"testing"
)

func TestChatEnabled(t *testing.T) {
	tests := map[string]struct {
		value    string
		setEnv   bool
		expected bool
	}{
		"enabled returns true": {
			value:    "enabled",
			setEnv:   true,
			expected: true,
		},
		"other value returns false": {
			value:    "disabled",
			setEnv:   true,
			expected: false,
		},
		"missing value returns false": {
			setEnv:   false,
			expected: false,
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			if testData.setEnv {
				t.Setenv(chatEnv, testData.value)
			} else {
				originalValue, hadOriginalValue := os.LookupEnv("CHAT")
				if err := os.Unsetenv(chatEnv); err != nil {
					t.Fatalf("error unsetting %s: %v", chatEnv, err)
				}
				t.Cleanup(func() {
					if hadOriginalValue {
						_ = os.Setenv(chatEnv, originalValue)
					}
				})
			}

			actual := ChatEnabled(nil)
			if actual != testData.expected {
				t.Errorf("ChatEnabled() returned %v, expected %v", actual, testData.expected)
			}
		})
	}
}
