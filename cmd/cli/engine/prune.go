package engine

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

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

func (cmd *pruneCommand) run(_ *cobra.Command, _ []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return err
	}
	activeEngineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
	if err != nil {
		return err
	}

	var componentsWithEnginesToRemove map[string][]string
	var componentsToRemove []string

	switch {
	case cmd.engine == "":
		componentsWithEnginesToRemove, err = cmd.getAllComponentsToRemove(*activeEngineManifest)
		if err != nil {
			return err
		}
		if !cmd.printComponentsAndConfirm(componentsWithEnginesToRemove, false) {
			return nil
		}
		return cmd.pruneAllInactiveEngines(slices.Collect(maps.Keys(componentsWithEnginesToRemove)))

	case cmd.engine == activeEngine:
		return fmt.Errorf("cannot prune the active engine '%s'", activeEngine)

	default:
		engineManifest, err := engines.LoadManifest(cmd.EnginesDir, cmd.engine)
		if err != nil {
			if errors.Is(err, engines.ErrManifestNotFound) {
				if cmd.Verbose {
					fmt.Println(err)
				}
				return fmt.Errorf("%q not found", cmd.engine)
			}
			return fmt.Errorf("error loading engine manifest: %v", err)
		}

		componentsWithEnginesToRemove, err = cmd.calculateRemovableComponents([]engines.Manifest{*engineManifest}, *activeEngineManifest)
		if err != nil {
			return err
		}
		if !cmd.printComponentsAndConfirm(componentsWithEnginesToRemove, true) {
			return nil
		}
		componentsToRemove = slices.Collect(maps.Keys(componentsWithEnginesToRemove))
		return cmd.pruneEngine(componentsToRemove, *engineManifest)
	}
}

func (cmd *pruneCommand) calculateRemovableComponents(enginesToCheck []engines.Manifest, activeEngineManifest engines.Manifest) (map[string][]string, error) {
	componentsEnginesMap := make(map[string][]string)

	activeSet := make(map[string]bool, len(activeEngineManifest.Components))
	for _, c := range activeEngineManifest.Components {
		activeSet[c] = true
	}
	for _, eng := range enginesToCheck {
		if eng.Name == activeEngineManifest.Name {
			continue
		}
		for _, component := range eng.Components {
			if activeSet[component] {
				continue
			}
			installed, err := componentInstalled(component)
			if err != nil {
				return nil, err
			}
			if installed {
				componentsEnginesMap[component] = append(componentsEnginesMap[component], eng.Name)
			}
		}
	}
	return componentsEnginesMap, nil
}

func (cmd *pruneCommand) getAllComponentsToRemove(activeEngineManifest engines.Manifest) (map[string][]string, error) {
	enginesToCheck, err := engines.LoadManifests(cmd.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifests: %w", err)
	}
	return cmd.calculateRemovableComponents(enginesToCheck, activeEngineManifest)
}

func (cmd *pruneCommand) pruneEngine(componentsToRemove []string, engine engines.Manifest) error {
	if err := unsetEngineConfig(engine.Name, cmd.Context); err != nil {
		return err
	}

	installed := make([]string, 0, len(componentsToRemove))
	for _, component := range componentsToRemove {
		if ok, err := componentInstalled(component); err == nil && ok {
			installed = append(installed, component)
		}
	}

	if len(installed) != 0 {
		if err := snapctl.RemoveComponents(installed...).Run(); err != nil {
			return fmt.Errorf("failed to remove components: %w", err)
		}
	}
	return nil
}

func (cmd *pruneCommand) pruneAllInactiveEngines(componentsToRemove []string) error {
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return err
	}
	var allEngines []engines.Manifest
	allEngines, err = engines.LoadManifests(cmd.EnginesDir)
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

func (cmd *pruneCommand) printComponentsAndConfirm(componentsWithEngines map[string][]string, isSingleEngine bool) bool {
	if len(componentsWithEngines) == 0 {
		fmt.Println("No components to remove.")
		return false
	}

	fmt.Println("Removing components:")

	componentSizes, err := snap_store.ComponentSizes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: unable to get component sizes: %v\n", err)
	}

	enginesSet := make(map[string]struct{})
	for componentName, engineNames := range componentsWithEngines {
		componentLine := componentName
		if size, ok := componentSizes[componentName]; ok {
			componentLine += fmt.Sprintf(" (%s)", utils.FmtBytes(uint64(size)))
		}

		if isSingleEngine {
			fmt.Printf("- %s\n", componentLine)
			continue
		}

		fmt.Printf("- %s [%s]\n", componentLine, strings.Join(engineNames, ", "))
		for _, engineName := range engineNames {
			enginesSet[engineName] = struct{}{}
		}
	}

	var confirmationPromptSentence string
	if isSingleEngine {
		confirmationPromptSentence = fmt.Sprintf("Continue pruning %v engine?", cmd.engine)
	} else {
		confirmationPromptSentence = fmt.Sprintf("Continue pruning %v engines?", slices.Collect(maps.Keys(enginesSet)))
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println()
		if !common.ConfirmationPrompt(confirmationPromptSentence) {
			fmt.Println("Exiting. No changes applied.")
			return false
		}
	}

	return true
}
