package commands

import (
	"fmt"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
)

func TestInfoLong(t *testing.T) {
	engine, err := engines.LoadManifest("../../../test_data/engines", "intel-gpu")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	cmd := showEngineCommand{
		format: "yaml",
	}
	err = cmd.printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoShort(t *testing.T) {
	engine, err := engines.LoadManifest("../../../test_data/engines", "cpu-avx1")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	cmd := showEngineCommand{
		format: "yaml",
	}
	err = cmd.printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func Example_showEngineCommand_printEngineManifestYaml() {
	engineManifest, err := engines.LoadManifest("../../../test_data/engines", "cuda-generic")
	if err != nil {
		panic(fmt.Sprintf("failed to load engine manifest: %v", err))
	}
	info, err := hardware_info.GetFromRawData("dummy-machine", true, "../../../test_data")
	if err != nil {
		panic(fmt.Sprintf("failed to get hardware info: %v", err))
	}
	scoredEngines, err := selector.ScoreEngines(info, []engines.Manifest{*engineManifest})
	if err != nil {
		panic(fmt.Sprintf("failed to score engines: %v", err))
	}

	cmd := showEngineCommand{format: "yaml"}
	if err := cmd.printEngineManifest(scoredEngines[0]); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// name: cuda-generic
	// description: Nvidia GPUs using CUDA. All major CUDA versions are targeted.
	// vendor: Canonical Ltd
	// grade: stable
	// devices:
	//     allof:
	//         - type: cpu
	//           architecture: amd64
	//           flags:
	//             - sse4_2
	//             - f16c
	//             - fma
	//             - avx
	//             - avx2
	//           compatibility-issues:
	//             - flag sse4_2 missing
	//             - flag f16c missing
	//             - flag fma missing
	//             - flag avx missing
	//             - flag avx2 missing
	//         - type: gpu
	//           bus: pci
	//           vendor-id: "0x10DE"
	//           vram: 5G
	//           compatibility-issues:
	//             - device not found
	// memory: 2G
	// disk-space: 5G
	// components:
	//     - dummy-component-1
	//     - dummy-component-2
	// configurations: {}
	// score: 0
	// compatible: false
	// compatibility-issues:
	//     - required device not found
}

func Example_showEngineCommand_printEngineManifestJson() {
	engineManifest, err := engines.LoadManifest("../../../test_data/engines", "cuda-generic")
	if err != nil {
		panic(fmt.Sprintf("failed to load engine manifest: %v", err))
	}
	info, err := hardware_info.GetFromRawData("dummy-machine", true, "../../../test_data")
	if err != nil {
		panic(fmt.Sprintf("failed to get hardware info: %v", err))
	}
	scoredEngines, err := selector.ScoreEngines(info, []engines.Manifest{*engineManifest})
	if err != nil {
		panic(fmt.Sprintf("failed to score engines: %v", err))
	}

	cmd := showEngineCommand{format: "json"}
	if err := cmd.printEngineManifest(scoredEngines[0]); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// {
	//   "name": "cuda-generic",
	//   "description": "Nvidia GPUs using CUDA. All major CUDA versions are targeted.",
	//   "vendor": "Canonical Ltd",
	//   "grade": "stable",
	//   "devices": {
	//     "anyof": null,
	//     "allof": [
	//       {
	//         "type": "cpu",
	//         "architecture": "amd64",
	//         "flags": [
	//           "sse4_2",
	//           "f16c",
	//           "fma",
	//           "avx",
	//           "avx2"
	//         ],
	//         "compatibility-issues": [
	//           "flag sse4_2 missing",
	//           "flag f16c missing",
	//           "flag fma missing",
	//           "flag avx missing",
	//           "flag avx2 missing"
	//         ]
	//       },
	//       {
	//         "type": "gpu",
	//         "bus": "pci",
	//         "vendor-id": "0x10DE",
	//         "vram": "5G",
	//         "compatibility-issues": [
	//           "device not found"
	//         ]
	//       }
	//     ]
	//   },
	//   "memory": "2G",
	//   "disk-space": "5G",
	//   "components": [
	//     "dummy-component-1",
	//     "dummy-component-2"
	//   ],
	//   "configurations": null,
	//   "score": 0,
	//   "compatible": false,
	//   "compatibility-issues": [
	//     "required device not found"
	//   ]
	// }
}
