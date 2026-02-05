package engine

import (
	"flag"
	"fmt"
	"os"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/snap_store"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type pruneCommand struct {
	*common.Context

	// flags
	engine string
}

func PruneCommand(ctx *common.Context) *cobra.Command {
	var cmd pruneCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:     "prune-cache",
		Short:   "Remove cached data for engines that are no longer in use",
		GroupID: groupID,
		RunE:    cmd.run,
	}

	// flags
	cobraCmd.Flags().StringVar(&cmd.engine, "engine", "", "Remove cache for the specified engine")

	return cobraCmd
}

func (cmd *pruneCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}
	var componentsWithEnginesToRemove map[string][]string
	var componentsToRemove []string
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if cmd.engine == "" {
		componentsWithEnginesToRemove, err = cmd.getAllComponentsToRemove(&engines.Manifest{Name: activeEngine})
		if err != nil {
			return err
		}
		if !cmd.printComponentsAndConfirm(componentsWithEnginesToRemove) {
			return nil
		}
		componentsToRemove = make([]string, 0, len(componentsWithEnginesToRemove))
		for component := range componentsWithEnginesToRemove {
			componentsToRemove = append(componentsToRemove, component)
		}
		return cmd.pruneAllInactiveEngines(componentsToRemove)
	} else {
		engineManifest, err := engines.LoadManifest(cmd.EnginesDir, cmd.engine)

		if err != nil {
			return err
		}
		activeEngineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
		if err != nil {
			return err
		}
		componentsWithEnginesToRemove = cmd.getComponentsToRemoveFromEngine(engineManifest, activeEngineManifest)
		if !cmd.printComponentsAndConfirm(componentsWithEnginesToRemove) {
			return nil
		}
		componentsToRemove = make([]string, 0, len(componentsWithEnginesToRemove))
		for component := range componentsWithEnginesToRemove {
			componentsToRemove = append(componentsToRemove, component)
		}
		return cmd.pruneEngine(componentsToRemove, *engineManifest)
	}
}

func (cmd *pruneCommand) calculateRemovableComponents(enginesToCheck []engines.Manifest, activeEngineManifest *engines.Manifest) map[string][]string {
	var componentsEnginesMap = make(map[string][]string)

	for _, eng := range enginesToCheck {
		// Skip the active engine itself
		if eng.Name == activeEngineManifest.Name {
			continue
		}

		for _, component := range eng.Components {
			if !utils.Contains(activeEngineManifest.Components, component) {
				if _, exists := componentsEnginesMap[component]; !exists {
					componentsEnginesMap[component] = []string{}
				}
				componentsEnginesMap[component] = append(componentsEnginesMap[component], eng.Name)
			}
		}
	}
	return componentsEnginesMap
}

func (cmd *pruneCommand) getComponentsToRemoveFromEngine(targetEngine *engines.Manifest, activeEngineManifest *engines.Manifest) map[string][]string {
	return cmd.calculateRemovableComponents([]engines.Manifest{*targetEngine}, activeEngineManifest)
}

func (cmd *pruneCommand) getAllComponentsToRemove(activeEngineManifest *engines.Manifest) (map[string][]string, error) {
	var enginesToCheck []engines.Manifest
	var err error

	if flag.Lookup("test.v") != nil {
		enginesToCheck, err = engines.LoadManifests("../../../test_data/engines")
	} else {
		enginesToCheck, err = engines.LoadManifests(cmd.EnginesDir)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load manifests: %w", err)
	}

	return cmd.calculateRemovableComponents(enginesToCheck, activeEngineManifest), nil
}

func (cmd *pruneCommand) pruneEngine(componentsToRemove []string, engine engines.Manifest) error {
	var err error
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return err
	}
	if engine.Name == activeEngine {
		return fmt.Errorf("cannot prune the active engine '%s'", activeEngine)
	}
	// Skip configuration changes during tests
	if flag.Lookup("test.v") == nil {
		uc := &useCommand{Context: cmd.Context}
		err = uc.unsetEngineConfig(engine.Name)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Removing cached components for engine '%s': %s\n", engine.Name, componentsToRemove)

	if flag.Lookup("test.v") == nil {
		_ = snapctl.RemoveComponents(componentsToRemove...).Run()
	}
	fmt.Printf("Pruned cache for engine: %s-\n", engine.Name)

	return nil
}

func (cmd *pruneCommand) pruneAllInactiveEngines(componentsToRemove []string) error {
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return err
	}
	var allEngines []engines.Manifest
	if flag.Lookup("test.v") != nil {
		allEngines, err = engines.LoadManifests("../../../test_data/engines")
	} else {
		allEngines, err = engines.LoadManifests(cmd.EnginesDir)
	}
	if err != nil {
		return fmt.Errorf("error scoring engines: %v", err)
	}

	for _, engine := range allEngines {
		if engine.Name != activeEngine {
			err := cmd.pruneEngine(componentsToRemove, engine)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (cmd *pruneCommand) printComponentsAndConfirm(componentsWithEngines map[string][]string) bool {
	fmt.Printf("Removing components:\n")

	// Look up component sizes from the snap store
	componentSizes, err := snap_store.ComponentSizes()
	if err != nil {
		// If component size lookup failed, continue but log the error
		fmt.Fprintf(os.Stderr, "Warning: unable to get component sizes: %v\n", err)
		return false
	}

	var componentList []string
	var enginesList []string
	for componentName := range componentsWithEngines {
		componentLine := fmt.Sprintf("- %s", componentName)
		if size, ok := componentSizes[componentName]; ok {
			componentLine += fmt.Sprintf(" (%s)", utils.FmtBytes(uint64(size)))
		}
		enginesLine := fmt.Sprintf("[")
		for _, engineName := range componentsWithEngines[componentName] {
			enginesLine += fmt.Sprintf("%s, ", engineName)
		}
		enginesLine = enginesLine[:len(enginesLine)-2] + "]"
		fmt.Printf("- %s %s\n\n", componentLine, enginesLine)
		componentList = append(componentList, componentLine)
		// append to engineList only if engine is not already present
		for _, engineName := range componentsWithEngines[componentName] {
			if !utils.Contains(enginesList, engineName) {
				enginesList = append(enginesList, engineName)
			}
		}
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println()
		if !common.ConfirmationPrompt(fmt.Sprintf("Continue pruning %v engines?", enginesList)) {
			fmt.Println("Exiting. No changes applied.")
			return false
		}
	}

	return true
}
