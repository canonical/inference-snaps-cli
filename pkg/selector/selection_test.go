package selector

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"gopkg.in/yaml.v3"
)

// Test that the expected engine is chosen from a list of engines
// device, engines, expected engine

type machineTopEngine struct {
	machine   string
	engines   []string
	topEngine string
}

var topEngineSets = []machineTopEngine{
	{
		// Ampere One machine should use ampere engine, not generic arm neon engine
		machine:   "ampere-one-m-banshee-12",
		engines:   []string{"ampere", "arm-neon"},
		topEngine: "ampere",
	},
	{
		// Ampere Altra machine should use ampere-altra engine, not generic arm neon engine
		machine:   "hp-proliant-rl300-gen11-altra",
		engines:   []string{"ampere-altra", "arm-neon"},
		topEngine: "ampere-altra",
	},
	{
		// Old CPU with Intel dGPU and NVIDIA dGPU - using intel-gpu because nvidia requires newer CPU flags
		machine:   "i5-3570k+arc-a580+gtx1080ti",
		engines:   []string{"cpu-avx1", "cuda-generic", "intel-cpu", "intel-gpu"},
		topEngine: "intel-gpu",
	},
	{
		// Machine with Intel CPU and Intel GPU should use GPU
		machine:   "mustang",
		engines:   []string{"cpu-avx1", "cpu-avx2", "cpu-avx512", "intel-cpu", "intel-gpu"},
		topEngine: "intel-gpu",
	},
	{
		// Machine with Intel iGPU and NVIDIA dGPU - always try and offload to dGPU if possible
		machine:   "system76-addw4",
		engines:   []string{"cpu-avx1", "cpu-avx2", "cuda-generic", "intel-cpu", "intel-gpu"},
		topEngine: "cuda-generic",
	},
	{
		// Machine with avx2 should prefer avx2 engine
		machine:   "xps13-7390",
		engines:   []string{"cpu-avx1", "cpu-avx2"},
		topEngine: "cpu-avx2",
	},
	{
		// Machine with Intel CPU should prefer intel-cpu engine above generic cpu engines
		machine:   "xps13-7390",
		engines:   []string{"cpu-avx1", "cpu-avx2", "intel-cpu"},
		topEngine: "intel-cpu",
	},
}

func TestTopEngine(t *testing.T) {
	for _, testSet := range topEngineSets {
		t.Run(testSet.machine+"/"+testSet.topEngine, func(t *testing.T) {
			var manifests []engines.Manifest
			for _, engineName := range testSet.engines {
				manifestFile := fmt.Sprintf("../../test_data/engines/%s/engine.yaml", engineName)
				data, err := os.ReadFile(manifestFile)
				if err != nil {
					t.Fatal(err)
				}

				var manifest engines.Manifest
				err = yaml.Unmarshal(data, &manifest)
				if err != nil {
					t.Fatal(err)
				}

				manifests = append(manifests, manifest)
			}

			hardwareInfo, err := hardware_info.GetFromRawData(t, testSet.machine, true)
			if err != nil {
				t.Fatal(err)
			}

			scoredEngines, err := ScoreEngines(hardwareInfo, manifests)
			if err != nil {
				t.Fatal(err)
			}

			topEngine, err := TopEngine(scoredEngines)
			if err != nil {
				t.Fatal(err)
			}

			if topEngine.Name != testSet.topEngine {
				for _, engine := range scoredEngines {
					t.Logf("%s=%d %s", engine.Name, engine.Score, strings.Join(engine.Notes, ", "))
				}
				t.Errorf("Top engine name: %s, expected: %s", topEngine.Name, testSet.topEngine)
			}
		})
	}
}

func TestLoadManifestFromDir_InvalidEngine(t *testing.T) {
	// Test loading an engine that doesn't exist
	_, err := LoadManifestFromDir("../../test_data/engines", "nonexistent-engine")
	if err == nil {
		t.Fatal("Expected error when loading non-existent engine, got nil")
	}
	
	// Verify the error message is user-friendly and doesn't expose file paths
	expectedMsg := `engine "nonexistent-engine" not found`
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
	
	// Ensure the error message doesn't contain file path details
	if strings.Contains(err.Error(), "/engine.yaml") || strings.Contains(err.Error(), "open") {
		t.Errorf("Error message should not expose file paths, got: %q", err.Error())
	}
}
