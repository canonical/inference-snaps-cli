package commands

import (
	"fmt"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func prepareTestData() (*listEnginesCommand, *outputEngines, error) {
	cache := storage.NewMockCache()
	err := cache.SetActiveEngine("example-memory")
	if err != nil {
		return nil, nil, fmt.Errorf("Error setting active engine name: %v", err)
	}

	allEngines, err := engines.LoadManifests("../../../test_data/engines")
	if err != nil {
		return nil, nil, fmt.Errorf("error loading engines: %v", err)
	}

	hardwareInfo, err := hardware_info.GetFromRawData("xps13-7390", true, "../../../test_data")
	if err != nil {
		return nil, nil, fmt.Errorf("error getting hardware info: %v", err)
	}

	scoredEngines, err := selector.ScoreEngines(hardwareInfo, allEngines)
	if err != nil {
		return nil, nil, fmt.Errorf("error scoring engines: %v", err)
	}

	// cmd.printEnginesTable needs to call `cmd.Cache.GetActiveEngine()` to get the current active engine
	// We therefore need to pass in the cache as context to `cmd`
	ctx := &common.Context{
		EnginesDir: "",
		Cache:      cache,
		Config:     nil,
	}
	cmd := listEnginesCommand{Context: ctx}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %v", common.LookingUpActiveEngine, err)
	}

	enginesList := outputEngines{
		ActiveEngine: activeEngine,
	}

	for _, se := range scoredEngines {
		enginesList.Engines = append(enginesList.Engines, common.NewEngineDetails(se))
	}

	return &cmd, &enginesList, nil
}

func TestList(t *testing.T) {
	cmd, enginesList, err := prepareTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	err = cmd.printEnginesJson(*enginesList)
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.printEnginesTable(*enginesList)
	if err != nil {
		t.Fatal(err)
	}
}
