package commands

import (
	"errors"
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/models"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type useModelCommand struct {
	*common.Context

	// flags
	auto      bool
	fix       bool
	assumeYes bool
	noRestart bool
}

func UseModel(ctx *common.Context) *cobra.Command {
	var cmd useModelCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "use-model [<model>]",
		Short: "Select a model",
		// Args
		// cli use-engine <engine> requires 1 argument
		// cli use-engine --auto does not support any arguments
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.auto, "auto", false, "automatically select the default model")
	cobraCmd.Flags().BoolVar(&cmd.fix, "fix", false, "fix issues with the currently active model")
	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for all prompts")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after changing model")

	return cobraCmd
}

func (cmd *useModelCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	// TODO
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (cmd *useModelCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	if cmd.auto {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both model name and --auto flag")
		}
		return cmd.autoSelectModel()
	} else if cmd.fix {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both model name and --fix flag")
		}
		// If no model is active, there's nothing to fix.
		_, err := common.FixActiveModel(cmd.Context)
		if errors.Is(err, common.ErrNoActiveModel) {
			return nil
		}
		return err
	} else {
		if len(args) == 1 {
			return cmd.switchModel(args[0])
		} else {
			return fmt.Errorf("model name not specified")
		}
	}
}

func (cmd *useModelCommand) switchModel(modelName string) error {

	availableModels, err := models.AvailableModelNames(cmd.ModelsDir)
	if err != nil {
		return fmt.Errorf("%s: %w", "loading available models", err)
	}
	if !slices.Contains(availableModels, modelName) {
		return fmt.Errorf("model %s does not exist", modelName)
	}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
	if err != nil {
		return fmt.Errorf("%s: %w", "loading engine manifest", err)
	}
	supportedModels := engineManifest.Model.Options

	if !slices.Contains(supportedModels, modelName) {
		return fmt.Errorf("model %s not supported by engine %s", modelName, activeEngine)
	}

	cancelledByUser, err := common.InstallMissingComponents(cmd.Context, cmd.assumeYes, engineManifest, modelName)
	if err != nil {
		return fmt.Errorf("installing missing components: %v", err)
	}

	if cancelledByUser {
		return nil
	}

	activeModelName, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("%s: %w", "looking up active model", err)
	}

	if activeModelName == modelName {
		// Engine not changed, nothing left to do
		return nil
	}

	if err = cmd.Cache.SetActiveModel(modelName); err != nil {
		return fmt.Errorf("setting active model: %v", err)
	}

	fmt.Printf("Model changed to %q.\n", modelName)

	// Ask if the user wants to restart
	if !cmd.noRestart {
		return common.PromptRestartToApplyChanges(cmd.Context, cmd.assumeYes)
	}

	return nil
}

func (cmd *useModelCommand) autoSelectModel() error {
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
	if err != nil {
		return fmt.Errorf("%s: %w", common.LoadingEngineManifest, err)
	}

	err = cmd.Cache.SetActiveModel(engineManifest.Model.Default)
	if err != nil {
		return fmt.Errorf("setting active model: %v", err)
	}

	// TODO check if the default model size will fit, otherwise check the next smaller one, iteratively

	fmt.Println("Switching to default model", engineManifest.Model.Default, "for engine", activeEngine)

	return cmd.switchModel(engineManifest.Model.Default)
}
