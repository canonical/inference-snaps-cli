package engine

import (
	"errors"
	"fmt"
	"os"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/spf13/cobra"
)

const groupID = "engine"

func Group(title string) *cobra.Group {
	return &cobra.Group{
		ID:    groupID,
		Title: title,
	}
}

func scoreEngines(ctx *common.Context) ([]engines.ScoredManifest, error) {
	allEngines, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("error loading engines: %v", err)
	}

	machineInfo, err := ctx.Cache.GetMachineInfo()
	if err != nil {
		return nil, fmt.Errorf("error getting machine info: %v", err)
	}

	// score engines
	scoredEngines, err := selector.ScoreEngines(machineInfo, allEngines)
	if err != nil {
		return nil, fmt.Errorf("error scoring engines: %v", err)
	}

	return scoredEngines, nil
}

func componentInstalled(component string) (bool, error) {
	// Check in /snap/$SNAP_INSTANCE_NAME/components/$SNAP_REVISION if component is mounted
	directoryPath := fmt.Sprintf("/snap/%s/components/%s/%s", env.SnapInstanceName(), env.SnapRevision(), component)

	info, err := os.Stat(directoryPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("error checking component directory %q: %v", component, err)
		}
	} else {
		if info.IsDir() {
			return true, nil
		} else {
			return false, fmt.Errorf("component %q exists but is not a directory", component)
		}
	}
}

func unsetEngineConfig(engineName string, context *common.Context) error {
	// Unset all engine configurations
	err := context.Config.Unset(".", storage.EngineConfig)
	if err != nil {
		return fmt.Errorf("error un-setting engine configurations: %v", err)
	}

	engine, err := engines.LoadManifest(context.EnginesDir, engineName)
	if err != nil {
		if errors.Is(err, engines.ErrManifestNotFound) {
			// TODO: remove this when implementing per-engine configuration
			// We can't know what user overrides were set if the manifest is missing
			fmt.Fprintf(os.Stderr, "Warning: previously active engine %q not found; skipping user configuration cleanup.\n", engineName)
			return nil
		}
		return fmt.Errorf("error loading engine manifest: %v", err)
	} else {
		// Unset any user overrides
		for k := range engine.Configurations {
			err = context.Config.Unset(k, storage.UserConfig)
			if err != nil {
				return fmt.Errorf("error un-setting configuration %q: %v", k, err)
			}
		}
	}

	return nil
}
