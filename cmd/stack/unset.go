package main

import (
	"fmt"

	"github.com/canonical/go-snapctl"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "unset <key>",
		Short: "Unset configuration option",
		// Long:  "",
		Args: cobra.ExactArgs(1),
		RunE: unset,
	}
	rootCmd.AddCommand(cmd)
}

func unset(_ *cobra.Command, args []string) error {
	return unsetValue(args[0])
}

func unsetValue(key string) error {
	err := snapctl.Unset(key).Run()
	if err != nil {
		return fmt.Errorf("error unsetting value of '%s': %v", key, err)
	}

	return nil
}
