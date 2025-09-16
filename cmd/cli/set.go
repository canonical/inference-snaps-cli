package main

import (
	"fmt"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/utils"
	"github.com/spf13/cobra"
)

func addSetCommand() {
	cmd := &cobra.Command{
		Use:   "set <key>",
		Short: "Set configuration option",
		// Long:  "",
		GroupID:           "config",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              set,
	}
	rootCmd.AddCommand(cmd)
}

func set(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return fmt.Errorf("No permission to set configurations. Please try again using sudo.")
	}
	return setValue(args[0])
}

func setValue(keyValue string) error {
	if keyValue[0] == '=' {
		return fmt.Errorf("key must not start with an equal sign")
	}

	// The value itself can contain an equal sign, so we split only on the first occurrence
	parts := strings.SplitN(keyValue, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("expected key=value, got %q", keyValue)
	}
	key, value := parts[0], parts[1]

	err := snapctl.Set(key, value).Run()
	if err != nil {
		return fmt.Errorf("error setting value %q for %q: %v", value, key, err)
	}

	return nil
}
