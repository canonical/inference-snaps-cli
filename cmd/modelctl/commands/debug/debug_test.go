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
		args             []string
		expectErr        bool
		expectedInErr    string
		expectedInOutput string
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
			args:             []string{"debug"},
			expectErr:        false,
			expectedInOutput: "Developer/debugging commands",
		},
		"help flag": {
			args:             []string{"debug", "--help"},
			expectErr:        false,
			expectedInOutput: "Developer/debugging commands",
		},
		// The parent now carries cobra.NoArgs and a RunE, so these two cases
		// guard the other direction: a real subcommand must still be reached,
		// and its own argument validation must be what runs, not the parent's.
		"valid subcommand": {
			args:             []string{"debug", "lint-package", "--help"},
			expectErr:        false,
			expectedInOutput: "Validate engines, models, and runtimes manifest files",
		},
		"valid subcommand keeps its own argument validation": {
			args:          []string{"debug", "lint-package"},
			expectErr:     true,
			expectedInErr: "accepts 1 arg(s), received 0",
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
			if !strings.Contains(output.String(), testData.expectedInOutput) {
				t.Errorf("Execute(%v) printed %q, expected it to contain %q",
					testData.args, output.String(), testData.expectedInOutput)
			}
		})
	}
}
