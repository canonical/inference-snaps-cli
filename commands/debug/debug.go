package debug

import (
	icmd "github.com/canonical/inference-snaps-cli/v2/internal/commands"
	"github.com/spf13/cobra"
)

func DebugCommand(ctx *icmd.Context) *cobra.Command {
	debugCmd := &cobra.Command{
		Use:    "debug",
		Long:   "Developer/debugging commands",
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
