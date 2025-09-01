package main

import (
	"log"

	"github.com/canonical/go-snapctl/env"
	"github.com/spf13/cobra"
)

var (
	stacksDir = env.Snap() + "/stacks"
	// rootCmd is the base command
	// It gets populated with subcommands via init functions
	rootCmd = &cobra.Command{
		Use:          env.SnapInstanceName(),
		SilenceUsage: true,
	}
)

func main() {
	cobra.EnableCommandSorting = false

	rootCmd.AddGroup(&cobra.Group{ID: "basics", Title: "Basic Commands:"})
	addStatusCommand()
	addChatCommand()

	rootCmd.AddGroup(&cobra.Group{ID: "config", Title: "Configuration Commands:"})
	addGetCommand()
	addSetCommand()
	addUnsetCommand()

	rootCmd.AddGroup(&cobra.Group{ID: "variants", Title: "Management Commands:"})
	addListCommand()
	addInfoCommand()
	addUseCommand()
	addPruneCommand()

	// disable logging timestamps
	log.SetFlags(0)

	// set a dummy root command if not in a snap
	if rootCmd.Use == "" {
		rootCmd.Use = "app"
	}

	// Hide the 'completion' command from help text
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.Execute()
}
