package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/validate"
	"github.com/spf13/cobra"
)

func init() {
	debugCmd := &cobra.Command{
		Use:   "debug",
		Short: "Debugging commands",
		// Long:  "",
		// GroupID: "debug",
	}

	machineInfoCmd := &cobra.Command{
		Use:               "machine-info",
		Short:             "Print machine information",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              machineInfo,
	}
	debugCmd.AddCommand(machineInfoCmd)

	validateCmd := &cobra.Command{
		Use:   "validate-variants",
		Short: "Validate variant manifest files",
		Args:  cobra.MinimumNArgs(1),
		RunE:  validateStackManifests,
	}
	debugCmd.AddCommand(validateCmd)

	rootCmd.AddCommand(debugCmd)
}

func machineInfo(_ *cobra.Command, args []string) error {
	hwInfo, err := hardware_info.Get(true)
	if err != nil {
		return fmt.Errorf("failed to get hardware info: %s", err)
	}

	jsonString, err := json.MarshalIndent(hwInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %s", err)
	}

	// print the JSON output
	fmt.Println(string(jsonString))

	return nil
}

func validateStackManifests(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no variant manifest specified")
	}

	for _, manifestPath := range args {
		err := validate.Stack(manifestPath)
		if err != nil {
			fmt.Printf("❌ %s: %s\n", manifestPath, err)
		} else {
			fmt.Printf("✅ %s\n", manifestPath)
		}
	}

	return nil
}
