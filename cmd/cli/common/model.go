package common

import (
	"fmt"
	"os"
	"slices"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
)

// FixActiveModel checks if the currently selected model is supported by the currently selected engine
// If it is not supported, the model is changed to the default one specified by the engine
// TODO if switching from llamacpp to ovms, and we have an e4b-gguf model selected, change to the ovms e4b-int4 model, base don the name field in the model manifest
// Returns the ID of the newly selected model
func FixActiveModel(ctx *Context) (string, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return "", fmt.Errorf("getting active engine: %v", err)
	}

	if activeEngineName == "" {
		// No engine is active, so there is no model to fix
		return "", nil
	}

	engineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return "", fmt.Errorf("loading active engine manifest: %v", err)
	}

	activeModelId, err := ctx.Cache.GetActiveModel()
	if err != nil {
		return "", ErrNoActiveModel
	}

	if activeModelId == "" {
		// No model is active, using default
		activeModelId = engineManifest.Model.Default
		err = ctx.Cache.SetActiveModel(activeModelId)
		if err != nil {
			return "", fmt.Errorf("setting active model: %v", err)
		}
	}

	// If the active model is unsupported by the engine, switch to the default one
	if !slices.Contains(engineManifest.Model.Options, activeModelId) {
		fmt.Fprintf(os.Stderr, "Warning: active model %q is not supported by engine %q, switching to default model %q\n", activeModelId, engineManifest.Name, engineManifest.Model.Default)
		activeModelId = engineManifest.Model.Default
		err = ctx.Cache.SetActiveModel(activeModelId)
		if err != nil {
			return "", fmt.Errorf("setting active model: %v", err)
		}
	}

	return activeModelId, nil
}
