package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func addRunCommand() {
	cmd := &cobra.Command{
		Use:    "run <path>",
		Short:  "Run a subprocess",
		Hidden: true,
		Args:   cobra.MaximumNArgs(1),
		RunE:   run,
	}
	rootCmd.AddCommand(cmd)
}

func run(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("unexpected number of arguments, expected 1 got %d", len(args))
	}

	// TODO
	// --wait-for-components

	err := loadEngineEnvironment()
	if err != nil {
		return fmt.Errorf("error loading engine environment: %v", err)
	}

	path := args[0]

	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

		err := godotenv.Overload(componentEnvFile)
		if err != nil {
			return fmt.Errorf("error loading env file for component %q: %v", componentName, err)
		}
	}

	return nil
}
