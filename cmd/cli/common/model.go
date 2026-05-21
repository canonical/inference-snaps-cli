package common

import (
	"fmt"
	"os"
	"slices"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
)

func FixActiveModel(ctx *Context) (string, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return "", fmt.Errorf("getting active engine: %v", err)
	}

	engineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return "", fmt.Errorf("loading active engine manifest: %v", err)
	}

	activeModel, err := ctx.Cache.GetActiveModel()
	if err != nil {
		return "", ErrNoActiveModel
	}

	if activeModel == "" {
		// No model is active, setting default
		activeModel = engineManifest.Model.Default
		err = ctx.Cache.SetActiveModel(activeModel)
		if err != nil {
			return "", fmt.Errorf("setting active model: %v", err)
		}
	}

	// If the active model is unsupported by the engine, switch to the default one
	if !slices.Contains(engineManifest.Model.Options, activeModel) {
		fmt.Fprintf(os.Stderr, "Warning: active model %q is not supported by engine %q, switching to default model %q\n", activeModel, engineManifest.Name, engineManifest.Model.Default)
		activeModel = engineManifest.Model.Default
		err = ctx.Cache.SetActiveModel(activeModel)
		if err != nil {
			return "", fmt.Errorf("setting active model: %v", err)
		}
	}

	return activeModel, nil
}
