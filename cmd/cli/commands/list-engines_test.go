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

func Example_printEnginesJson() {
	cmd, enginesList, err := prepareTestData()
	if err != nil {
		panic(fmt.Sprintf("Error preparing test data: %v", err))
	}

	err = cmd.printEnginesJson(*enginesList)
	if err != nil {
		panic(fmt.Sprintf("Error printing engines json: %v", err))
	}

	// Output:
	// {
	//   "active-engine": "example-memory",
	//   "engines": [
	//     {
	//       "name": "amd-gpu",
	//       "description": "AMD specific engine targeting only one microarchitecture.",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64"
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x1002",
	//             "microarchitecture": "gfx1032",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           }
	//         ]
	//       },
	//       "memory": "2G",
	//       "disk-space": "5G",
	//       "components": [
	//         "dummy-component-1",
	//         "dummy-component-2"
	//       ],
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "ampere",
	//       "description": "Test ampere selection",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "arm64",
	//             "implementer-id": "0xC0",
	//             "compatibility-issues": [
	//               "architecture not arm64"
	//             ]
	//           }
	//         ]
	//       },
	//       "memory": "4G",
	//       "disk-space": "5G",
	//       "components": null,
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "ampere-altra",
	//       "description": "Test ampere selection",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "arm64",
	//             "implementer-id": "0x41",
	//             "part-number": "0xD0C",
	//             "compatibility-issues": [
	//               "architecture not arm64"
	//             ]
	//           },
	//           {
	//             "bus": "pci",
	//             "vendor-id": "0x1DEF",
	//             "device-id": "0xE100",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           }
	//         ]
	//       },
	//       "memory": "4G",
	//       "disk-space": "5G",
	//       "components": null,
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "arm-neon",
	//       "description": "ARM CPUs with NEON instruction set",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "arm64",
	//             "features": [
	//               "asimd"
	//             ],
	//             "compatibility-issues": [
	//               "architecture not arm64"
	//             ]
	//           }
	//         ],
	//         "allof": null
	//       },
	//       "memory": "2G",
	//       "disk-space": "4G",
	//       "components": null,
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "cpu-avx1",
	//       "description": "Legacy CPUs with only SSE4.2 (2008) and AVX (2011)",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "flags": [
	//               "sse4_2",
	//               "avx"
	//             ]
	//           }
	//         ],
	//         "allof": null
	//       },
	//       "memory": "4G",
	//       "disk-space": "8G",
	//       "components": null,
	//       "configurations": null,
	//       "score": 14,
	//       "compatible": true
	//     },
	//     {
	//       "name": "cpu-avx2",
	//       "description": "CPUs with AVX2",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "flags": [
	//               "sse4_2",
	//               "f16c",
	//               "fma",
	//               "avx",
	//               "avx2"
	//             ]
	//           }
	//         ]
	//       },
	//       "memory": "8G",
	//       "disk-space": "5G",
	//       "components": null,
	//       "configurations": null,
	//       "score": 17,
	//       "compatible": true
	//     },
	//     {
	//       "name": "cpu-avx512",
	//       "description": "CPUs with AVX512",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "flags": [
	//               "sse4_2",
	//               "f16c",
	//               "fma",
	//               "avx",
	//               "avx2",
	//               "avx512f",
	//               "avx512dq",
	//               "avx512bw"
	//             ],
	//             "compatibility-issues": [
	//               "flag avx512f missing",
	//               "flag avx512dq missing",
	//               "flag avx512bw missing"
	//             ]
	//           }
	//         ]
	//       },
	//       "memory": "8G",
	//       "disk-space": "5G",
	//       "components": null,
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "cpu-devel",
	//       "description": "Requires any CPU but is grade devel",
	//       "vendor": "Canonical Ltd",
	//       "grade": "devel",
	//       "devices": {
	//         "anyof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64"
	//           },
	//           {
	//             "type": "cpu",
	//             "architecture": "arm64",
	//             "compatibility-issues": [
	//               "architecture not arm64"
	//             ]
	//           }
	//         ],
	//         "allof": null
	//       },
	//       "memory": "1G",
	//       "disk-space": "2G",
	//       "components": [
	//         "dummy-component-1",
	//         "dummy-component-2",
	//         "dummy-component-3"
	//       ],
	//       "configurations": null,
	//       "score": 12,
	//       "compatible": true
	//     },
	//     {
	//       "name": "cuda-generic",
	//       "description": "Nvidia GPUs using CUDA. All major CUDA versions are targeted.",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "flags": [
	//               "sse4_2",
	//               "f16c",
	//               "fma",
	//               "avx",
	//               "avx2"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "bus": "pci",
	//             "vendor-id": "0x10DE",
	//             "vram": "5G",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           }
	//         ]
	//       },
	//       "memory": "2G",
	//       "disk-space": "5G",
	//       "components": [
	//         "dummy-component-1",
	//         "dummy-component-2"
	//       ],
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "example-memory",
	//       "description": "Legacy CPUs, offering full accuracy but very high memory usage",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "manufacturer-id": "AuthenticAMD",
	//             "compatibility-issues": [
	//               "manufacturer id mismatch: GenuineIntel"
	//             ]
	//           },
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "manufacturer-id": "GenuineIntel"
	//           }
	//         ],
	//         "allof": null
	//       },
	//       "memory": "35G",
	//       "disk-space": "29G",
	//       "components": [
	//         "dummy-component-3"
	//       ],
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "insufficient memory"
	//       ]
	//     },
	//     {
	//       "name": "intel-cpu",
	//       "description": "Use Intel CPUs",
	//       "vendor": "Intel Corporation",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "manufacturer-id": "GenuineIntel"
	//           }
	//         ]
	//       },
	//       "memory": "5G",
	//       "disk-space": "10G",
	//       "components": [
	//         "dummy-component-1"
	//       ],
	//       "configurations": null,
	//       "score": 18,
	//       "compatible": true
	//     },
	//     {
	//       "name": "intel-gpu",
	//       "description": "Modern Intel GPUs (\u003e=gen 13)",
	//       "vendor": "Intel Corporation",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": [
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F80",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F81",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F82",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F83",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F84",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F85",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F86",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F87",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x4F88",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5690",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5691",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5692",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5693",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5694",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5695",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5696",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5697",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x5698",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A0",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A1",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A2",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A3",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A4",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A5",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A6",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A7",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A8",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56A9",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56B0",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56B1",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56B2",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56B3",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56BA",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56BB",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56BC",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56BD",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56BE",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x56BF",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xE202",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xE20B",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xE20C",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xE20D",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xE210",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xE212",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xE215",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xE216",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA720",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA721",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA780",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA781",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA782",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA783",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA788",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA789",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA78A",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA78B",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA7A0",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA7A1",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA7A8",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA7A9",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA7AA",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA7AB",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA7AC",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xA7AD",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7DD5",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7D40",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7D55",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7D60",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7D45",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x6420",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x64A0",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x64B0",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB640",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7D41",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7D51",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7D67",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7DD1",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB080",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB081",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB082",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB083",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB08F",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB090",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB0A0",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x8086",
	//             "device-id": "0xB0B0",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           }
	//         ],
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64"
	//           }
	//         ]
	//       },
	//       "memory": "4G",
	//       "disk-space": "8G",
	//       "components": null,
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "intel-npu",
	//       "description": "Intel NPUs",
	//       "vendor": "Intel Corporation",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": [
	//           {
	//             "bus": "pci",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x643E",
	//             "snap-connections": [
	//               "intel-npu",
	//               "npu-libs"
	//             ],
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "bus": "pci",
	//             "vendor-id": "0x8086",
	//             "device-id": "0x7D1D",
	//             "snap-connections": [
	//               "intel-npu",
	//               "npu-libs"
	//             ],
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           }
	//         ],
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64"
	//           }
	//         ]
	//       },
	//       "memory": "4G",
	//       "disk-space": "8G",
	//       "components": null,
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "rocm-generic",
	//       "description": "AMD GPUs using ROCm. All major ROCm versions are targeted.",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": [
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x1002",
	//             "microarchitecture": "gfx1152",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x1002",
	//             "microarchitecture": "gfx1030",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           }
	//         ],
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64"
	//           }
	//         ]
	//       },
	//       "memory": "2G",
	//       "disk-space": "5G",
	//       "components": [
	//         "dummy-component-1",
	//         "dummy-component-2"
	//       ],
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     }
	//   ]
	// }
}
