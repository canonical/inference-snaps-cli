package debug

import (
	"bytes"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/spf13/cobra"
)

// newTestRoot mirrors the production command tree: the debug command is always
// attached to the root command. The topology matters, because Cobra's default
// argument validator only reports unknown subcommands for a command without a
// parent, so a standalone debug command would pass these tests regardless.
func newTestRoot() *cobra.Command {
	rootCmd := &cobra.Command{Use: "cli"}
	rootCmd.AddCommand(DebugCommand(&common.Context{}))

	return rootCmd
}

func TestDebugSubcommandValidation(t *testing.T) {
	tests := map[string]struct {
		args          []string
		expectErr     bool
		expectedInErr string
	}{
		"unknown subcommand": {
			args:          []string{"debug", "validate-engines"},
			expectErr:     true,
			expectedInErr: `unknown command "validate-engines" for "cli debug"`,
		},
		"unknown subcommand with arguments": {
			args:          []string{"debug", "validate-engines", "."},
			expectErr:     true,
			expectedInErr: `unknown command "validate-engines" for "cli debug"`,
		},
		"no subcommand prints help": {
			args:      []string{"debug"},
			expectErr: false,
		},
		"help flag": {
			args:      []string{"debug", "--help"},
			expectErr: false,
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			rootCmd := newTestRoot()

			var output bytes.Buffer
			rootCmd.SetOut(&output)
			rootCmd.SetErr(&output)
			rootCmd.SetArgs(testData.args)

			err := rootCmd.Execute()

			if testData.expectErr {
				if err == nil {
					t.Fatalf("Execute(%v) returned no error, expected one", testData.args)
				}
				if !strings.Contains(err.Error(), testData.expectedInErr) {
					t.Errorf("Execute(%v) returned error %q, expected it to contain %q",
						testData.args, err.Error(), testData.expectedInErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Execute(%v) returned error %v, expected none", testData.args, err)
			}
			if !strings.Contains(output.String(), "Developer/debugging commands") {
				t.Errorf("Execute(%v) did not print the command help", testData.args)
			}
		})
	}
}
