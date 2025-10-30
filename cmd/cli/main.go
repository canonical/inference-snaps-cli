package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	enginesDir       = env.Snap() + "/engines"
	snapInstanceName = env.SnapInstanceName()
	// rootCmd is the base command
	// It gets populated with subcommands
	rootCmd = &cobra.Command{
		Use:          snapInstanceName,
		SilenceUsage: true,
		Long:         "", // Base command description TBA
	}

	cache  = storage.NewCache()
	config = storage.NewConfig()

	// Error types
	ErrPermissionDenied = errors.New("permission denied, try again with sudo")
)

func main() {
	cobra.EnableCommandSorting = false

	// TODO: refact: functions called below add to the global rootCmd

	rootCmd.AddGroup(&cobra.Group{ID: "basics", Title: "Basic Commands:"})
	addStatusCommand()
	addChatCommand()

	rootCmd.AddGroup(&cobra.Group{ID: "config", Title: "Configuration Commands:"})
	addGetCommand()
	addSetCommand()

	rootCmd.AddGroup(&cobra.Group{ID: "engines", Title: "Management Commands:"})
	addListCommand()
	addInfoCommand()
	addUseCommand()

	// other commands (help is added by default)
	addShowMachineCommand()
	addDebugCommand()

	// disable logging timestamps
	log.SetFlags(0)

	// set a dummy root command if not in a snap
	if rootCmd.Use == "" {
		rootCmd.Use = "cli"
	}

	// Hide the 'completion' command from help text
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func loadEngineEnvironment() error {
	activeEngineName, err := cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("error looking up active engine: %v", err)
	}

	if activeEngineName == "" {
		return nil
	}

	manifest, err := selector.LoadManifestFromDir(enginesDir, activeEngineName)
	if err != nil {
		return fmt.Errorf("error loading active engine manifest: %v", err)
	}

	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return fmt.Errorf("SNAP_COMPONENTS environment variable not set")
	}

	const envFile = "component.env"
	for _, componentName := range manifest.Components {
		componentEnvFile := filepath.Join(componentsDir, componentName, envFile)

		err := godotenv.Load(componentEnvFile)
		if err != nil {
			return fmt.Errorf("error loading env file for component %q: %v", componentName, err)
		}
	}

}
