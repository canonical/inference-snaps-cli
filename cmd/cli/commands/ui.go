package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

type uiCommand struct {
	*common.Context
}

func Ui(ctx *common.Context) *cobra.Command {
	var cmd uiCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "ui",
		Short:             "Launch web UI",
		Long:              "Open the snap's builtin web user interface in the default browser",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobraCmd
}

func (cmd *uiCommand) run(_ *cobra.Command, _ []string) error {
	// Having all the components installed is not required, but it provides good feedback
	if err := common.WaitForComponents(cmd.Context); err != nil {
		return fmt.Errorf("waiting for component: %s", err)
	}

	// Get web ui url
	url, err := common.UiServerHttpUrl(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting ui server url: %s", err)
	}

	// Check ui server
	services, err := common.ServiceStatuses()
	if err != nil {
		return fmt.Errorf("getting service statuses: %v", err)
	}
	serverStatus, ok := services["ui-server"]
	if !ok {
		return fmt.Errorf("ui-server: service not found")
	}

	if serverStatus != "active" {
		fmt.Fprintf(os.Stderr, "Warning: ui-server: service is not \"active\" (current status: %s)\n\n", serverStatus)
		fmt.Printf("Make sure the ui server is running, then open %s in your browser.\n", url)
		return nil
	}

	// xdg open
	err = exec.Command("xdg-open", url).Start()
	if err != nil {
		return fmt.Errorf("xdg-open: %v", err)
	}

	return nil
}
