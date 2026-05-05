package commands

import (
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func ExampleUseEngine_noRestartWhenEngineUnchanged() {
	cache := storage.NewMockCache()
	cache.SetActiveEngine("intel-gpu")
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		Context: &common.Context{
			EnginesDir: "../../../test_data/engines",
			Cache:      cache,
			Config:     config,
			Snap:       snap.Mock(),
		},
	}

	if err := cmd.switchEngine("intel-gpu"); err != nil {
		panic(err)
	}

	// Output:
}

func ExampleUseEngine_restartWhenEngineChanged() {
	cache := storage.NewMockCache()
	cache.SetActiveEngine("intel-gpu")
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		Context: &common.Context{
			EnginesDir: "../../../test_data/engines",
			Cache:      cache,
			Config:     config,
			Snap:       snap.Mock(),
		},
	}

	if err := cmd.switchEngine("cpu-avx1"); err != nil {
		panic(err)
	}

	// Output:
	// Engine changed to "cpu-avx1".
	// [mock] Restarting all services
}

func ExampleUseEngine_autoSelectEngine() {
	cache := storage.NewMockCache()
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		Context: &common.Context{
			EnginesDir: "../../../test_data/engines",
			Cache:      cache,
			Config:     config,
			Snap:       snap.Mock(),
		},
	}
	cmd.Cache.SetActiveEngine("")
	cmd.Verbose = true
	notCompatibleManifest, err := engines.LoadManifest(cmd.Context.EnginesDir, "not-compatible-engine")
	if err != nil {
		panic(err)
	}
	cpuDevelManifest, err := engines.LoadManifest(cmd.Context.EnginesDir, "cpu-devel")
	if err != nil {
		panic(err)
	}
	cpuManifest, err := engines.LoadManifest(cmd.Context.EnginesDir, "cpu")
	if err != nil {
		panic(err)
	}

	allEngines := []engines.Manifest{*notCompatibleManifest, *cpuDevelManifest, *cpuManifest}
	machineInfo, err := hardware_info.GetFromRawData("mustang", true, "../../../test_data")
	if err != nil {
		panic(err)
	}

	scoredEngines, err := selector.ScoreEngines(machineInfo, allEngines)
	if err != nil {
		panic(err)
	}
	if err := cmd.autoSelectEngineWithScoredEngines(scoredEngines); err != nil {
		panic(err)
	}

	// Output:
	// Evaluating engines for optimal hardware compatibility:
	// ✘ not-compatible-engine: not compatible
	//   - required device not found
	// • cpu-devel: devel, score=12
	// ✔ cpu: compatible, score=12
	// Selected engine: cpu
	// Engine changed to "cpu".
	// [mock] Restarting all services
}
