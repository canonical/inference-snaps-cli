package debug

import (
	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/spf13/cobra"
)

func DebugCommand(ctx *common.Context) *cobra.Command {
	debugCmd := &cobra.Command{
		Use:  "debug",
		Long: "Developer/debugging commands",
		// Reject anything that is not one of the subcommands below. Cobra's
		// default validator only reports unknown subcommands for the root
		// command, and it skips argument validation entirely for commands
		// without a Run function, so both are needed here to avoid exiting 0
		// on an unrecognized subcommand.
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Hidden: true,
	}

	debugCmd.AddCommand(
		ValidateCommand(ctx),
		SelectCommand(ctx),
		ChatCommand(ctx),
		ServeWebUiCommand(ctx),
	)

	return debugCmd
}
