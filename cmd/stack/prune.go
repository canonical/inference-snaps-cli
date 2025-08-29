package main

import (
	"fmt"

	"github.com/canonical/go-snapctl"
	"github.com/spf13/cobra"
)

var (
	pruneAll bool
)

func init() {
	cmd := &cobra.Command{
		Use:   "prune-variant [<variant>]",
		Short: "Remove unused variant resources",
		// Long:  "",
		GroupID:           "variants",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: pruneValidArgs,
		RunE:              prune,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&pruneAll, "all", false, "remove resources of all unused variants")

	rootCmd.AddCommand(cmd)
}

func pruneValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	stacksJson, err := snapctl.Get("stacks").Document().Run()
	if err != nil {
		fmt.Printf("Error loading variants: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	stacks, err := parseStacksJson(stacksJson)
	if err != nil {
		fmt.Printf("Error parsing variants: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	usedStack, err := snapctl.Get("stack").Run()
	if err != nil {
		fmt.Printf("Error loading used variant: %v", err)
	}

	var stackNames []cobra.Completion
	for i := range stacks {
		// Can prune any stack except the currently used one
		if stacks[i].Name != usedStack {
			stackNames = append(stackNames, stacks[i].Name)
		}
	}
	if len(stackNames) == 0 {
		// No stacks available to prune
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return stackNames, cobra.ShellCompDirectiveNoFileComp
}

func prune(_ *cobra.Command, args []string) error {
	return fmt.Errorf("Pruning has not been implemented.\n")
}
