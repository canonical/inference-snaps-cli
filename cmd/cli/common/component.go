package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
)

func ComponentInstalled(component string) (bool, error) {
	// Check in /snap/$SNAP_INSTANCE_NAME/components/$SNAP_REVISION if component is mounted
	directoryPath := fmt.Sprintf("/snap/%s/components/%s/%s", env.SnapInstanceName(), env.SnapRevision(), component)

	info, err := os.Stat(directoryPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("checking component directory %q: %v", component, err)
		}
	} else {
		if info.IsDir() {
			return true, nil
		} else {
			return false, fmt.Errorf("component %q exists but is not a directory", component)
		}
	}
}

func WaitForComponents(ctx *Context) error {
	const maxWait = 3600 // seconds
	const interval = 10  // seconds
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}

	if activeEngineName == "" {
		return ErrNoActiveEngine
	}

	manifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return fmt.Errorf("loading engine manifest: %v", err)
	}

	missing, err := checkMissingComponents(manifest)
	if err != nil {
		return err
	}

	for elapsed := 0; elapsed < maxWait && len(missing) > 0; elapsed += interval {
		fmt.Fprintf(os.Stderr, "Waiting for required snap components: %s (%d/%ds)\n",
			strings.Join(missing, ", "), elapsed, maxWait)

		time.Sleep(interval * time.Second)

		missing, err = checkMissingComponents(manifest)
		if err != nil {
			return err
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("timeout after waiting %ds for required components: %s",
			maxWait, strings.Join(missing, ", "))
	}

	return nil
}

// TODO: unify with similar code in use-engine.go
func checkMissingComponents(manifest *engines.Manifest) ([]string, error) {
	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return nil, fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	var missing []string
	for _, component := range manifest.Components {
		componentPath := filepath.Join(componentsDir, component)
		if _, err := os.Stat(componentPath); os.IsNotExist(err) {
			missing = append(missing, component)
		}
	}

	return missing, nil
}

func ProcessPassthroughConfigs(ctx *Context) error {
	configs, err := ctx.Config.GetAll()
	if err != nil {
		return fmt.Errorf("getting passthrough configs: %v", err)
	}

	// Filter configs for passthrough keys starting with "passthrough." But remove the "passthrough." prefix for easier processing by the engine components.
	passthroughConfigs := make(map[string]any)
	for k, v := range configs {
		if strings.HasPrefix(k, "passthrough.") {
			passthroughConfigs[strings.TrimPrefix(k, "passthrough.")] = v
		}
	}

	// Set environment variables for passthrough configs. The engine components will read these environment variables.
	for k, v := range passthroughConfigs {
		if strings.HasPrefix(k, "environment.") {
			envVarName := strings.TrimPrefix(k, "environment.")
			envVarValue := fmt.Sprintf("%v", v)
			// envVarName is kebab-case and lower case, environment variables conventionally are upper case and snake case, so convert kebab-case to snake_case and uppercase the env var name
			envVarName = strings.ToUpper(strings.ReplaceAll(envVarName, "-", "_"))
			err := os.Setenv(envVarName, envVarValue)
			if err != nil {
				return fmt.Errorf("setting environment variable for passthrough config %q: %v", k, err)
			}
		}
	}
	return nil

}
