package main

import (
	"fmt"

	"github.com/canonical/go-snapctl"
	"github.com/spf13/cobra"
)

var getCmd *cobra.Command

func init() {
	getCmd = &cobra.Command{
		Use:   "get <key>",
		Short: "Print configuration option",
		// Long:  "Print configuration option",
		Args: cobra.ExactArgs(1),
		RunE: get,
	}
	rootCmd.AddCommand(getCmd)
}

func get(cmd *cobra.Command, args []string) error {
	return getValue(args[0])
}

func getValue(key string) error {
	value, err := snapctl.Get(key).Run()
	if err != nil {
		return fmt.Errorf("error getting value of '%s': %v", key, err)
	}

	if value == "" {
		return fmt.Errorf("no value set for key '%s'", key)
	}

	// print config value
	fmt.Println(value)

	return nil
}
