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

func TestGetEnginesTable(t *testing.T) {
	cmd, enginesList, err := prepareTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	tableStr, err := cmd.getEnginesTable(*enginesList)
	if err != nil {
		t.Fatalf("Error getting engines table: %v", err)
	}

	expectedTable := `ENGINE           VENDOR             DESCRIPTION                          COMPAT
intel-cpu        Intel Corporation  Use Intel CPUs                       yes   
cpu-avx2         Canonical Ltd      CPUs with AVX2                       yes   
cpu-avx1         Canonical Ltd      Legacy CPUs with only SSE4.2 (2008…  yes   
cpu-devel        Canonical Ltd      Requires any CPU but is grade devel  devel 
intel-gpu        Intel Corporation  Modern Intel GPUs (>=gen 13)         no    
example-memory*  Canonical Ltd      Legacy CPUs, offering full accurac…  no    
cuda-generic     Canonical Ltd      Nvidia GPUs using CUDA. All major …  no    
rocm-generic     Canonical Ltd      AMD GPUs using ROCm. All major ROC…  no    
arm-neon         Canonical Ltd      ARM CPUs with NEON instruction set   no    
ampere-altra     Canonical Ltd      Test ampere selection                no    
ampere           Canonical Ltd      Test ampere selection                no    
intel-npu        Intel Corporation  Intel NPUs                           no    
amd-gpu          Canonical Ltd      AMD specific engine targeting only…  no    
cpu-avx512       Canonical Ltd      CPUs with AVX512                     no    
`

	if *tableStr != expectedTable {
		t.Errorf("Engine table not as expected.\n\nGot:\n\n%s\n\nWant:\n\n%s", *tableStr, expectedTable)
	}
}
