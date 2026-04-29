package commands

import (
	"fmt"
	"strings"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type setCommand struct {
	*common.Context

	// flags
	packageConfig bool
	engineConfig  bool
	assumeYes     bool
	noRestart     bool
}

func Set(ctx *common.Context) *cobra.Command {
	var cmd setCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "set <key=value>...",
		Short:             "Set configurations",
		Long:              "Set a configuration",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.packageConfig, "package", false, "set package configurations")
	if err := cobraCmd.Flags().MarkHidden("package"); err != nil {
		panic(err)
	}
	cobraCmd.Flags().BoolVar(&cmd.engineConfig, "engine", false, "set engine configuration")
	if err := cobraCmd.Flags().MarkHidden("engine"); err != nil {
		panic(err)
	}
	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for all prompts")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after setting the configuration")

	return cobraCmd
}

func (cmd *setCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	return cmd.setValues(args)
}

func (cmd *setCommand) setValues(keyValues []string) error {

	currentValues := map[string]string{}
	currentKnown := map[string]bool{}
	seenKeys := map[string]bool{}
	finalValues := map[string]string{}

	// Validate key values
	for _, keyValue := range keyValues {
		key, value, err := cmd.parseKeyValue(keyValue)
		if err != nil {
			return err
		}

		// Reject duplicate keys
		if seenKeys[key] {
			return fmt.Errorf("duplicate key: %q", key)
		}
		seenKeys[key] = true

		currentValue, found, err := cmd.validateAndGetCurrentValue(key)
		if err != nil {
			return err
		}

		currentValues[key] = currentValue
		currentKnown[key] = found
		finalValues[key] = value
	}

	// Apply configurations
	anyChange := false
	for k, v := range finalValues {
		if err := cmd.setConfigs(k, v); err != nil {
			return err
		}
		if !currentKnown[k] || currentValues[k] != v {
			anyChange = true
		}
	}

	// Restart if configurations were changed
	if anyChange {
		if err := cmd.restartToApply(); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *setCommand) parseKeyValue(keyValue string) (key, value string, err error) {
	if keyValue == "" {
		return "", "", fmt.Errorf("expected key=value, got %q", keyValue)
	}

	if keyValue[0] == '=' {
		return "", "", fmt.Errorf("key must not start with an equal sign")
	}

	// The value itself can contain an equal sign, so we split only on the first occurrence
	parts := strings.SplitN(keyValue, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected key=value, got %q", keyValue)
	}
	return parts[0], parts[1], nil
}

func (cmd *setCommand) validateAndGetCurrentValue(key string) (string, bool, error) {
	currValMap, err := cmd.Config.Get(key)
	if err != nil {
		return "", false, fmt.Errorf("checking existing keys: %s", err)
	}
	currVal, found := currValMap[key]
	if !found && !strings.HasPrefix(key, "passthrough.") {
		return "", false, fmt.Errorf("unknown key: %q", key)
	}

	if !found {
		return "", false, nil
	}

	return fmt.Sprint(currVal), true, nil
}

func (cmd *setCommand) setConfigs(key, value string) error {
	var confType storage.ConfigType
	switch {
	case cmd.packageConfig:
		confType = storage.PackageConfig
	case cmd.engineConfig:
		confType = storage.EngineConfig
	default:
		confType = storage.UserConfig
	}

	if err := cmd.Config.Set(key, value, confType); err != nil {
		return fmt.Errorf("setting %q to %q: %v", key, value, err)
	}
	return nil
}

func (cmd *setCommand) restartToApply() error {
	if !cmd.noRestart {
		msg := fmt.Sprintf("Restart %s to apply the changes?", cmd.Snap.InstanceName())
		if cmd.assumeYes || common.PromptYN(msg, true) {
			if err := cmd.Snap.Restart(); err != nil {
				return fmt.Errorf("restarting snap: %v", err)
			}
		}
	}
	return nil
}
