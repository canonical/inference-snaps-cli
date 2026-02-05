package engine

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestPrune(t *testing.T) {
	var err error
	cache := storage.NewMockCache()
	err = cache.SetActiveEngine("example-memory")
	if err != nil {
		t.Fatalf("Error setting active engine name: %v", err)
	}

	allEngines, err := engines.LoadManifests("../../../test_data/engines")
	if err != nil {
		t.Fatalf("error loading engines: %v", err)
	}

	hardwareInfo, err := hardware_info.GetFromRawData(t, "xps13-7390", true, "../../../test_data")
	if err != nil {
		t.Fatalf("error getting hardware info: %v", err)
	}

	scoredEngines, err := selector.ScoreEngines(hardwareInfo, allEngines)
	if err != nil {
		t.Fatalf("error scoring engines: %v", err)
	}

	// cmd.printEnginesTable needs to call `cmd.Cache.GetActiveEngine()` to get the current active engine
	// We therefore need to pass in the cache as context to `cmd`
	ctx := &common.Context{
		EnginesDir: "",
		Cache:      cache,
		Config:     nil,
	}
	cmd := pruneCommand{Context: ctx}
	var activeEngineManifest engines.Manifest
	activeEngineName, err := cache.GetActiveEngine()
	for _, scoredEngine := range scoredEngines {
		if scoredEngine.Name == activeEngineName {
			activeEngineManifest = scoredEngine.Manifest
			break
		}
	}

	var componentsToRemoveWithEngines map[string][]string
	var componentsToRemove []string
	// test pruning active engine
	componentsToRemoveWithEngines = cmd.getComponentsToRemoveFromEngine(&activeEngineManifest, &activeEngineManifest)
	componentsToRemove = make([]string, 0, len(componentsToRemoveWithEngines))
	for component := range componentsToRemoveWithEngines {
		componentsToRemove = append(componentsToRemove, component)
	}
	err = cmd.pruneEngine(componentsToRemove, activeEngineManifest)
	if err == nil {
		t.Fatalf("error pruning ACTIVE engine: %v", err)
	}

	// test pruning inactive engine
	componentsToRemoveWithEngines = cmd.getComponentsToRemoveFromEngine(&scoredEngines[0].Manifest, &activeEngineManifest)
	componentsToRemove = make([]string, 0, len(componentsToRemoveWithEngines))
	for component := range componentsToRemoveWithEngines {
		componentsToRemove = append(componentsToRemove, component)
	}
	err = cmd.pruneEngine(componentsToRemove, scoredEngines[0].Manifest)
	if err != nil {
		t.Fatalf("error pruning engine: %v", err)
	}

	// test pruning all inactive engines
	componentsToRemoveWithEngines, err = cmd.getAllComponentsToRemove(&activeEngineManifest)
	if err != nil {
		t.Fatalf("error getting all components to remove: %v", err)
	}
	componentsToRemove = make([]string, 0, len(componentsToRemoveWithEngines))
	for component := range componentsToRemoveWithEngines {
		componentsToRemove = append(componentsToRemove, component)
	}
	err = cmd.pruneAllInactiveEngines(componentsToRemove)
	if err != nil {
		t.Fatalf("error pruning all inactive engines: %v", err)
	}
}
