package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/visualization"
	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/host"
	"github.com/spf13/cobra"
)

type showMachineCommand struct {
	*common.Context

	// flags
	format string
}

func ShowMachine(ctx *common.Context) *cobra.Command {
	var cmd showMachineCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "show-machine",
		Short:             "Print information about the host machine",
		Long:              "Print information about the host machine, including hardware and compute resources",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	supportedFormats := []string{"json", "yaml"}
	cobraCmd.Flags().StringVar(
		&cmd.format,
		"format",
		"yaml",
		fmt.Sprintf("output format (%s)", strings.Join(supportedFormats, ", ")),
	)

	return cobraCmd
}

func (cmd *showMachineCommand) run(_ *cobra.Command, _ []string) error {
	info, err := cmd.fetchMachineInfoWithSpinner()
	if err != nil {
		return err
	}

	return cmd.printMachineInfo(info)
}

func (cmd *showMachineCommand) printMachineInfo(info *machine.MachineInfo) error {
	rendered, err := visualization.Marshal(visualization.New(info), visualization.Format(cmd.format))
	if err != nil {
		return err
	}
	fmt.Printf("%s", rendered)
	return nil
}

func (cmd *showMachineCommand) fetchMachineInfoWithSpinner() (*machine.MachineInfo, error) {
	stopProgress := common.StartProgressSpinner("Gathering machine information")
	hwInfo, warnings, err := machine.Get(host.Real(), true)
	stopProgress()

	if len(warnings) > 0 && cmd.Verbose {
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("getting machine info: %s", err)
	}

	return hwInfo, nil
}
