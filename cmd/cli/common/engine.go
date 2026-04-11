package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"gopkg.in/yaml.v3"
)

const (
	componentEnv    = "COMPONENT"
	ProgressScoring = "Checking engines"
)

type ComponentLayout struct {
	Symlink string `yaml:"symlink"`
}

type ComponentSettings struct {
	componentName  string
	Servers        map[string]map[string]string `yaml:"servers"`
	Environment    []string                     `yaml:"environment"`
	Layout         map[string]ComponentLayout   `yaml:"layout"`
	expandedLayout map[string]ComponentLayout
}

func EngineComponentSettings(ctx *Context) ([]ComponentSettings, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}

	if activeEngineName == "" {
		return nil, ErrNoActiveEngine
	}

	manifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("loading engine manifest: %v", err)
	}

	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return nil, fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	var settingsCollection []ComponentSettings
	for _, componentName := range manifest.Components {
		componentPath := filepath.Join(componentsDir, componentName)
		componentYamlFile := filepath.Join(componentPath, "component.yaml")

		data, err := os.ReadFile(componentYamlFile)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", componentYamlFile, err)
		}

		var settings ComponentSettings
		err = yaml.Unmarshal(data, &settings)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling %s: %v", componentYamlFile, err)
		}

		settings.componentName = componentName

		settingsCollection = append(settingsCollection, settings)
	}

	return settingsCollection, nil
}

func loadEngineEnvironmentFromSettingsCollection(settingsCollection []ComponentSettings) error {
	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	for i, settings := range settingsCollection {
		// Set component path env var for expansion
		componentPath := filepath.Join(componentsDir, settings.componentName)
		if err := os.Setenv(componentEnv, componentPath); err != nil {
			return fmt.Errorf("setting env %q: %v", componentEnv, err)
		}

		for i := range settings.Environment {
			// Split into key/value
			kv := settings.Environment[i]
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid env var %q", kv)
			}
			k, v := parts[0], parts[1]

			// Expand all env vars in value
			v = os.ExpandEnv(v)

			err := os.Setenv(k, v)
			if err != nil {
				return fmt.Errorf("setting env var %q: %v", k, err)
			}
		}

		settingsCollection[i].expandedLayout = make(map[string]ComponentLayout, len(settings.Layout))
		for k, v := range settings.Layout {
			ComponentLayout := ComponentLayout{
				Symlink: os.ExpandEnv(v.Symlink),
			}
			settingsCollection[i].expandedLayout[os.ExpandEnv(k)] = ComponentLayout
		}

		for layoutPath, layout := range settingsCollection[i].expandedLayout {
			if layout.Symlink != "" {
				if err := createTemporarySymlink(layout.Symlink, layoutPath); err != nil {
					return fmt.Errorf("creating temporary symlink for component %q: %v", settings.componentName, err)
				}
			}
		}
	}

	if err := os.Unsetenv(componentEnv); err != nil {
		return fmt.Errorf("error unsetting %q: %v", componentEnv, err)
	}

	return nil
}

func pathWithinTmp(path string) (bool, error) {
	link, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("getting absolute path of %s: %v", path, err)
	}
	if strings.HasPrefix(link, "/tmp") {
		return true, nil
	}
	return false, nil
}

func removeTemporarySymlink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // If the file doesn't exist, consider it removed
		}
		return fmt.Errorf("stat %s: %v", path, err)
	}

	// Check if the path is within /tmp before removing
	if withinTmp, err := pathWithinTmp(path); err != nil {
		return fmt.Errorf("checking if path is within /tmp: %v", err)
	} else if !withinTmp {
		return fmt.Errorf("layout path outside of /tmp: %s", path)
	}

	// Only remove if it's a symlink to avoid accidentally deleting other files
	if info.Mode()&os.ModeSymlink != 0 {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("removing file %s: %v", path, err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Warning: expected %q to be a symlink but it is not; skipping removal.\n", path)
	}

	return nil
}

func createTemporarySymlink(target, link string) error {
	// Reject any layout path that is outside of /tmp
	if withinTmp, err := pathWithinTmp(link); err != nil {
		return fmt.Errorf("checking if path is within /tmp: %v", err)
	} else if !withinTmp {
		return fmt.Errorf("layout path outside of /tmp: %s", link)
	}

	// Create directory tree for the link
	if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
		return fmt.Errorf("creating directory for symlink %q: %v", link, err)
	}

	// Remove existing symlink if it exists
	if err := removeTemporarySymlink(link); err != nil {
		return fmt.Errorf("removing existing symlink at %q: %v", link, err)
	}

	// Create new symlink
	if err := os.Symlink(target, link); err != nil {
		return fmt.Errorf("creating symlink from %q to %q: %v", link, target, err)
	}

	return nil
}

func unloadEngineEnvironmentFromSettingsCollection(settingsCollection []ComponentSettings) error {

	// remove the symlinks created for the engine components
	for _, settings := range settingsCollection {
		for layoutPath := range settings.expandedLayout {
			if err := removeTemporarySymlink(layoutPath); err != nil {
				return fmt.Errorf("removing symlink %q: %v", layoutPath, err)
			}
		}
	}

	return nil
}

// LoadEngineEnvironment sets env vars of the active engine's components for the current process
// and creates any necessary symlinks
func LoadEngineEnvironment(ctx *Context) (func(), error) {
	settingsCollection, err := EngineComponentSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading engine component settings: %v", err)
	}
	err = loadEngineEnvironmentFromSettingsCollection(settingsCollection)
	return func() {
		if err := unloadEngineEnvironmentFromSettingsCollection(settingsCollection); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to unload engine environment: %v\n", err)
		}
	}, err
}

// SetEngineConfig sets configurations of the given engine.
// It does not unset previous engine configurations.
func SetEngineConfig(engine *engines.Manifest, ctx *Context) error {
	for confKey, confVal := range engine.Configurations {
		err := ctx.Config.SetDocument(confKey, confVal, storage.EngineConfig)
		if err != nil {
			return fmt.Errorf("setting engine configuration %q: %v", confKey, err)
		}
	}
	return nil
}

func UnsetEngineConfig(engineName string, unsetUserOverrides bool, ctx *Context) error {
	// Unset all engine configurations
	err := ctx.Config.Unset(".", storage.EngineConfig)
	if err != nil {
		return fmt.Errorf("un-setting engine configurations: %v", err)
	}

	if unsetUserOverrides {
		engine, err := engines.LoadManifest(ctx.EnginesDir, engineName)
		if err != nil {
			if errors.Is(err, engines.ErrManifestNotFound) {
				// TODO: remove this when implementing per-engine configuration
				// We can't know what user overrides were set if the manifest is missing
				fmt.Fprintf(os.Stderr, "Warning: previously active engine %q not found; skipping user configuration cleanup.\n", engineName)
				return nil
			}
			return fmt.Errorf("loading engine manifest: %v", err)
		} else {
			// Unset any user overrides
			for k := range engine.Configurations {
				err = ctx.Config.Unset(k, storage.UserConfig)
				if err != nil {
					return fmt.Errorf("un-setting configuration %q: %v", k, err)
				}
			}
		}
	}

	return nil
}

/*
ScoreEngines loads all engine manifests, looks up the host machine information,
and scores the engines according to their compatibility with the host.

Warning: calls to this function can block for a number of seconds while the host machine information is being looked up.
*/
func ScoreEngines(ctx *Context) ([]engines.ScoredManifest, error) {
	allEngines, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("loading engines: %v", err)
	}

	machineInfo, err := hardware_info.Get(false)
	if err != nil {
		return nil, fmt.Errorf("getting machine info: %v", err)
	}

	scoredEngines, err := selector.ScoreEngines(machineInfo, allEngines)
	if err != nil {
		return nil, fmt.Errorf("scoring engines: %v", err)
	}

	return scoredEngines, nil
}
