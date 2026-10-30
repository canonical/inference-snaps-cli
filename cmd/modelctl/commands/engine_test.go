package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

func TestEngineNoScoring(t *testing.T) {
	for _, newCmd := range []struct {
		name string
		new  func(*common.Context) *cobra.Command
	}{
		{"engine", Engine},
		{"show-engine", ShowEngine},
	} {
		for _, format := range []string{"yaml", "json"} {
			for _, active := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%s/active=%t", newCmd.name, format, active), func(t *testing.T) {
					cache := storage.NewMockCache()
					if err := cache.SetActiveEngine("cuda-generic"); err != nil {
						t.Fatal(err)
					}
					ctx := &common.Context{
						EnginesDir: "../../../test_data/engines",
						Cache:      cache,
					}
					cmd := newCmd.new(ctx)
					args := []string{"--no-scoring", "--format=" + format}
					if !active {
						args = append(args, "cuda-generic")
					}
					cmd.SetArgs(args)

					output, err := os.CreateTemp(t.TempDir(), "output")
					if err != nil {
						t.Fatal(err)
					}
					defer output.Close()
					originalStdout := os.Stdout
					os.Stdout = output
					t.Cleanup(func() { os.Stdout = originalStdout })
					if err := cmd.Execute(); err != nil {
						t.Fatal(err)
					}
					os.Stdout = originalStdout
					data, err := os.ReadFile(output.Name())
					if err != nil {
						t.Fatal(err)
					}

					manifest, err := engines.LoadManifest(ctx.EnginesDir, "cuda-generic")
					if err != nil {
						t.Fatal(err)
					}
					var expected []byte
					if format == "json" {
						expected, err = json.MarshalIndent(manifest, "", "  ")
						expected = append(expected, '\n')
					} else {
						expected, err = yaml.Marshal(manifest)
					}
					if err != nil {
						t.Fatal(err)
					}
					if string(data) != string(expected) {
						t.Fatalf("expected unscored manifest:\n%s\ngot:\n%s", expected, data)
					}
				})
			}
		}
	}
}

func TestEngineNoScoringErrors(t *testing.T) {
	t.Run("missing engine", func(t *testing.T) {
		cmd := engineCommand{
			Context:   &common.Context{EnginesDir: t.TempDir()},
			noScoring: true,
		}
		if err := cmd.engine("missing"); !errors.Is(err, engines.ErrManifestNotFound) {
			t.Fatalf("expected missing manifest error, got %v", err)
		}
	})
	t.Run("invalid manifest", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, "invalid"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "invalid", engines.ManifestFilename), []byte("name: ["), 0644); err != nil {
			t.Fatal(err)
		}
		cmd := engineCommand{
			Context:   &common.Context{EnginesDir: dir},
			noScoring: true,
		}
		if err := cmd.engine("invalid"); err == nil || !strings.Contains(err.Error(), "loading engine manifest") {
			t.Fatalf("expected manifest loading error, got %v", err)
		}
	})
	t.Run("no active engine", func(t *testing.T) {
		cmd := engineCommand{
			Context:   &common.Context{Cache: storage.NewMockCache()},
			noScoring: true,
		}
		if err := cmd.run(nil, nil); !errors.Is(err, common.ErrNoActiveEngine) {
			t.Fatalf("expected no active engine error, got %v", err)
		}
	})
	t.Run("unsupported format", func(t *testing.T) {
		cmd := engineCommand{
			Context:   &common.Context{EnginesDir: "../../../test_data/engines"},
			noScoring: true,
			format:    "invalid",
		}
		if err := cmd.engine("cuda-generic"); err == nil || !strings.Contains(err.Error(), "unknown format") {
			t.Fatalf("expected unsupported format error, got %v", err)
		}
	})
}

func TestInfoLong(t *testing.T) {
	engine, err := engines.LoadManifest("../../../test_data/engines", "intel-gpu")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	cmd := engineCommand{
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

	cmd := engineCommand{
		format: "yaml",
	}
	err = cmd.printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func scoreEngineAgainstMachine(engineName string, machineName string) (*engines.ScoredManifest, error) {
	engineManifest, err := engines.LoadManifest("../../../test_data/engines", engineName)
	if err != nil {
		return nil, fmt.Errorf("failed to load engine manifest: %v", err)
	}
	info, err := machineFixture(machineName)
	if err != nil {
		return nil, fmt.Errorf("failed to get hardware info: %v", err)
	}
	scoredEngines, err := selector.ScoreEngines(info, []engines.Manifest{*engineManifest})
	if err != nil {
		return nil, fmt.Errorf("failed to score engines: %v", err)
	}

	if len(scoredEngines) != 1 {
		return nil, fmt.Errorf("invalid scored engines count: %d, expected 1", len(scoredEngines))
	}

	return &scoredEngines[0], nil
}

func TestUnsupportedFormatResultsInError(t *testing.T) {
	engineManifest, err := scoreEngineAgainstMachine("cpu-avx1", "dummy-machine")
	if err != nil {
		t.Fatalf("could not score manifest: %v", err)
	}

	cmd := engineCommand{format: "invalid-format"}
	err = cmd.printEngineManifest(*engineManifest)

	if err == nil {
		t.Fatalf("expected unsupported format to error out, got nil error")
	}
}

func Example_engineCommand_printEngineManifestYaml() {
	engineManifest, err := scoreEngineAgainstMachine("cuda-generic", "dummy-machine")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := engineCommand{format: "yaml"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// name: cuda-generic
	// summary: CUDA generic engine
	// description: Nvidia GPUs using CUDA. All major CUDA versions are targeted.
	// vendor: Canonical Ltd
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
	// runtime: llama-cpp-cuda
	// model:
	//     default: 26b-q4-k-m-gguf
	//     options:
	//         - 26b-q4-k-m-gguf
	//         - 30b-a3b-q4-k-m-gguf
	// configurations:
	//     sleep-idle-seconds: 600
	// score: 0
	// compatible: false
	// compatibility-issues:
	//     - required device not found
}

func Example_engineCommand_printEngineManifestJson() {
	engineManifest, err := scoreEngineAgainstMachine("cuda-generic", "dummy-machine")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := engineCommand{format: "json"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// {
	//   "name": "cuda-generic",
	//   "summary": "CUDA generic engine",
	//   "description": "Nvidia GPUs using CUDA. All major CUDA versions are targeted.",
	//   "vendor": "Canonical Ltd",
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
	//   "runtime": "llama-cpp-cuda",
	//   "model": {
	//     "default": "26b-q4-k-m-gguf",
	//     "options": [
	//       "26b-q4-k-m-gguf",
	//       "30b-a3b-q4-k-m-gguf"
	//     ]
	//   },
	//   "configurations": {
	//     "sleep-idle-seconds": 600
	//   },
	//   "score": 0,
	//   "compatible": false,
	//   "compatibility-issues": [
	//     "required device not found"
	//   ]
	// }
}

func Example_engineCommand_printHappyEngineManifestYaml() {
	engineManifest, err := scoreEngineAgainstMachine("intel-cpu", "i7-1165G7")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := engineCommand{format: "yaml"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// name: intel-cpu
	// summary: Intel CPU engine
	// description: Use Intel CPUs
	// vendor: Intel Corporation
	// devices:
	//     allof:
	//         - type: cpu
	//           architecture: amd64
	//           manufacturer-id: GenuineIntel
	// runtime: openvino-model-server
	// model:
	//     default: 4b-it-int4-fq-ov
	//     options:
	//         - 4b-it-int4-fq-ov
	// configurations:
	//     target-device: CPU
	// score: 16
	// compatible: true
}

func Example_engineCommand_printHappyEngineManifestJson() {
	engineManifest, err := scoreEngineAgainstMachine("intel-cpu", "i7-1165G7")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := engineCommand{format: "json"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// {
	//   "name": "intel-cpu",
	//   "summary": "Intel CPU engine",
	//   "description": "Use Intel CPUs",
	//   "vendor": "Intel Corporation",
	//   "devices": {
	//     "anyof": null,
	//     "allof": [
	//       {
	//         "type": "cpu",
	//         "architecture": "amd64",
	//         "manufacturer-id": "GenuineIntel"
	//       }
	//     ]
	//   },
	//   "runtime": "openvino-model-server",
	//   "model": {
	//     "default": "4b-it-int4-fq-ov",
	//     "options": [
	//       "4b-it-int4-fq-ov"
	//     ]
	//   },
	//   "configurations": {
	//     "target-device": "CPU"
	//   },
	//   "score": 16,
	//   "compatible": true
	// }
}

func Example_engineCommand_printEngineManifestWithLineBreaksInDescriptionYaml() {
	engineManifest, err := scoreEngineAgainstMachine("amd-gpu", "dummy-machine")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := engineCommand{format: "yaml"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// name: amd-gpu
	// summary: AMD GPU engine
	// description: |
	//     AMD specific engine targeting one microarchitecture:
	//       - gfx1032
	// vendor: Canonical Ltd
	// devices:
	//     allof:
	//         - type: cpu
	//           architecture: amd64
	//         - type: gpu
	//           vendor-id: "0x1002"
	//           microarchitecture: gfx1032
	//           compatibility-issues:
	//             - device not found
	// configurations: {}
	// score: 0
	// compatible: false
	// compatibility-issues:
	//     - required device not found
}

func Example_engineCommand_printEngineManifestWithLineBreaksInDescriptionJson() {
	engineManifest, err := scoreEngineAgainstMachine("amd-gpu", "dummy-machine")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := engineCommand{format: "json"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// {
	//   "name": "amd-gpu",
	//   "summary": "AMD GPU engine",
	//   "description": "AMD specific engine targeting one microarchitecture:\n  - gfx1032\n",
	//   "vendor": "Canonical Ltd",
	//   "devices": {
	//     "anyof": null,
	//     "allof": [
	//       {
	//         "type": "cpu",
	//         "architecture": "amd64"
	//       },
	//       {
	//         "type": "gpu",
	//         "vendor-id": "0x1002",
	//         "microarchitecture": "gfx1032",
	//         "compatibility-issues": [
	//           "device not found"
	//         ]
	//       }
	//     ]
	//   },
	//   "runtime": "",
	//   "model": {
	//     "default": "",
	//     "options": null
	//   },
	//   "configurations": null,
	//   "score": 0,
	//   "compatible": false,
	//   "compatibility-issues": [
	//     "required device not found"
	//   ]
	// }
}
