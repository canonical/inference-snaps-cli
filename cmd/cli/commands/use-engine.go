package commands

import (
	"errors"
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type useEngineCommand struct {
	*common.Context

	// flags
	auto      bool
	fix       bool
	assumeYes bool
	noRestart bool
}

func UseEngine(ctx *common.Context) *cobra.Command {
	var cmd useEngineCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "use-engine [<engine>]",
		Short: "Select an engine",
		// Args
		// cli use-engine <engine> requires 1 argument
		// cli use-engine --auto does not support any arguments
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.auto, "auto", false, "automatically select a compatible engine")
	cobraCmd.Flags().BoolVar(&cmd.fix, "fix", false, "fix issues with the currently active engine")
	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for all prompts")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after changing engine")

	return cobraCmd
}

func (cmd *useEngineCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	manifests, err := engines.LoadManifests(cmd.EnginesDir)
	if err != nil {
		fmt.Printf("Error loading engines: %v\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var engineNames []cobra.Completion
	for i := range manifests {
		engineNames = append(engineNames, manifests[i].Name)
	}

	return engineNames, cobra.ShellCompDirectiveNoSpace
}

func (cmd *useEngineCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	if cmd.auto {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both engine name and --auto flag")
		}
		return cmd.autoSelectEngine()
	} else if cmd.fix {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both engine name and --fix flag")
		}
		// If no engine is active, there's nothing to fix.
		err := cmd.fixActiveEngine()
		if errors.Is(err, common.ErrNoActiveEngine) {
			return nil
		}
		return err
	} else {
		if len(args) == 1 {
			return cmd.switchEngine(args[0])
		} else {
			return fmt.Errorf("engine name not specified")
		}
	}
}

func (cmd *useEngineCommand) autoSelectEngine() error {
	scoredEngines, err := common.ScoreEnginesWithSpinner(cmd.Context)
	if err != nil {
		return fmt.Errorf("scoring engines: %v", err)
	}

	return cmd.autoSelectScoredEngine(scoredEngines)
}

func (cmd *useEngineCommand) autoSelectScoredEngine(scoredEngines []engines.ScoredManifest) error {

	fmt.Println("Evaluating engines for optimal hardware compatibility:")
	for _, engine := range scoredEngines {
		if engine.Score == 0 {
			fmt.Printf("✘ %s: not compatible\n", engine.Name)

			// Only print incompatibility reasons if verbose flag is set
			if cmd.Verbose {
				reasons := cmd.verboseIncompatibilityReasons(engine.CompatibilityReport)
				for _, reason := range reasons {
					fmt.Printf("  - %s\n", reason)
				}
			}
		} else if engine.Grade != "stable" {
			fmt.Printf("• %s: devel, score=%d\n", engine.Name, engine.Score)
		} else {
			fmt.Printf("✔ %s: compatible, score=%d\n", engine.Name, engine.Score)
		}
	}

	selectedEngine, err := selector.TopEngine(scoredEngines)
	if err != nil {
		return fmt.Errorf("finding top engine: %v", err)
	}

	fmt.Printf("Selected engine: %s\n", selectedEngine.Name)

	err = cmd.switchEngine(selectedEngine.Name)
	if err != nil {
		return fmt.Errorf("use engine: %s", err)
	}

	return nil
}

// switchEngine changes the engine that is used by the snap
func (cmd *useEngineCommand) switchEngine(engineName string) error {

	engine, err := engines.LoadManifest(cmd.EnginesDir, engineName)
	if err != nil {
		if errors.Is(err, engines.ErrManifestNotFound) {
			if cmd.Verbose {
				fmt.Println(err)
			}
			return fmt.Errorf("%q not found", engineName)
		}
		return fmt.Errorf("loading engine manifest: %v", err)
	}

	activeModel, err := common.FixActiveModel(cmd.Context)
	if err != nil {
		return err
	}

	cancelledByUser, err := common.InstallMissingComponents(cmd.Context, cmd.assumeYes, engine, activeModel)
	if err != nil {
		return fmt.Errorf("installing missing components: %v", err)
	}

	if cancelledByUser {
		return nil
	}

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	if activeEngineName == engineName {
		// Engine not changed, nothing left to do
		return nil
	}

	// Unset active engine's configurations
	if activeEngineName != "" {
		err = common.UnsetEngineConfig(activeEngineName, true, cmd.Context)
		if err != nil {
			return fmt.Errorf("un-setting engine configurations: %v", err)
		}
	}

	if err = cmd.Cache.SetActiveEngine(engine.Name); err != nil {
		return fmt.Errorf("setting active engine: %v", err)
	}

	if err = common.SetEngineConfig(engine, cmd.Context); err != nil {
		return fmt.Errorf("setting new engine configurations: %v", err)
	}

	fmt.Printf("Engine changed to %q.\n", engineName)

	// Currently we cannot reliably determine if the service is active to automatically restart it
	// See https://bugs.launchpad.net/snapd/+bug/2137543
	//
	// Ask if the user wants to restart
	if !cmd.noRestart {
		return common.PromptRestartToApplyChanges(cmd.Context, cmd.assumeYes)
	}

	return nil
}

func (cmd *useEngineCommand) fixActiveEngine() error {
	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}
	if activeEngineName == "" {
		return common.ErrNoActiveEngine
	}

	// If active engine no longer exist, auto select another one
	engine, err := engines.LoadManifest(cmd.EnginesDir, activeEngineName)
	if errors.Is(err, engines.ErrManifestNotFound) {
		fmt.Printf("Active engine %q not found, performing auto selection instead.\n", activeEngineName)
		return cmd.autoSelectEngine()
	} else if err != nil {
		return fmt.Errorf("loading active engine manifest: %v", err)
	}

	// Verify active model is supported or switch to default
	activeModel, err := common.FixActiveModel(cmd.Context)
	if err != nil {
		return err
	}

	// Make sure all components are correctly installed and engine is configured
	if _, err = common.InstallMissingComponents(cmd.Context, cmd.assumeYes, engine, activeModel); err != nil {
		return fmt.Errorf("installing missing components: %v", err)
	}

	if err = common.UnsetEngineConfig(activeEngineName, false, cmd.Context); err != nil {
		return fmt.Errorf("un-setting engine configurations: %v", err)
	}
	if err = common.SetEngineConfig(engine, cmd.Context); err != nil {
		return fmt.Errorf("setting engine configurations: %v", err)
	}

	return nil
}

func (cmd *useEngineCommand) verboseIncompatibilityReasons(report engines.CompatibilityReport) []string {
	var reasons []string
	if !report.CompatibleMemory {
		reasons = append(reasons, fmt.Sprintf("requires %s memory, has %s (RAM + swap)", utils.FmtBytes(report.RequiredMemory), utils.FmtBytes(report.TotalRAM+report.TotalSwap)))
	}
	if !report.CompatibleDisk {
		reasons = append(reasons, fmt.Sprintf("requires %s disk space, has %s", utils.FmtBytes(report.RequiredDiskSpace), utils.FmtBytes(report.AvailableDiskSpace)))
	}
	if !report.CompatibleDevices {
		reasons = append(reasons, "required device not found")
	}
	return reasons
}
