package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/cpu"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
	"github.com/canonical/lscompute/pkg/machine/disk"
	"github.com/canonical/lscompute/pkg/machine/memory"
)

func getTestMachine() *machine.Machine {
	return &machine.Machine{
		CPUs: []cpu.CPU{{
			Architecture:   "amd64",
			ManufacturerId: "GenuineIntel",
			Flags:          []string{"fpu", "vme", "de"},
		}},
		Memory: memory.Memory{TotalRam: 67012501504, TotalSwap: 0},
		Disk: []disk.Disk{{
			Total:     1006451294208,
			Available: 943543738368,
			Path:      "/var/lib/snapd/snaps",
		}},
		PCIDevices: []pci.Device{{
			Bus:                  "pci",
			Slot:                 "0000:00:00.0",
			BusNumber:            0x0,
			DeviceClass:          0x600,
			ProgrammingInterface: new(uint8(0)),
			VendorId:             0x8086,
			DeviceId:             0x4637,
			SubvendorId:          new(uint16(0x103C)),
			SubdeviceId:          new(uint16(0x89C6)),
			FriendlyNames: pci.FriendlyNames{
				VendorName:    "Intel Corporation",
				SubvendorName: "Hewlett-Packard Company",
			}},
		}}
}

// writeModelYAML creates a model manifest at modelsDir/<name>/model.yaml with the given content.
func writeModelYAML(t *testing.T, modelsDir, name, content string) {
	t.Helper()
	dir := filepath.Join(modelsDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "model.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile model.yaml: %v", err)
	}
}

// makeModelStatusCtx builds a Context with the given active model and a temp modelsDir.
func makeModelStatusCtx(t *testing.T, modelsDir, modelName string) *Context {
	t.Helper()
	cache := storage.NewMockCache()
	if err := cache.SetActiveModel(modelName); err != nil {
		t.Fatalf("SetActiveModel: %v", err)
	}
	return &Context{
		ModelsDir: modelsDir,
		Cache:     cache,
	}
}

func TestModelStatus_ModelNameWithEqualsInValue(t *testing.T) {
	modelsDir := t.TempDir()
	writeModelYAML(t, modelsDir, "my-model", `name: my=model`)
	ctx := makeModelStatusCtx(t, modelsDir, "my-model")

	status, err := ModelStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status["name"] != "my=model" {
		t.Errorf("expected name %q, got %q", "my=model", status["name"])
	}
}

func TestModelStatus_NonExistentModel(t *testing.T) {
	modelsDir := t.TempDir()
	ctx := makeModelStatusCtx(t, modelsDir, "non-existent")

	_, err := ModelStatus(ctx)
	if err == nil {
		t.Fatal("expected error for non-existent model, got nil")
	}
}

func TestGetModelManifestByNameOrAlias(t *testing.T) {
	tests := []struct {
		name         string
		activeEngine string
		modelYAML    string // empty means don't write a model manifest
		engineYAML   string // empty means don't write an engine manifest
		runtimeYAML  string // empty means don't write a runtime manifest
		query        string
		wantName     string // non-empty: expect this Name in the returned manifest
		wantAlias    string // non-empty: expect this Alias in the returned manifest
		wantErr      bool   // true: expect any non-nil error
	}{
		{
			name:         "found by alias",
			activeEngine: "my-engine",
			modelYAML:    "name: my-model-id\nalias: my-model\ndisk-size: 1G\n",
			engineYAML:   "name: my-engine\nruntime: my-runtime\nmodel:\n  options:\n    - my-model-id\n",
			runtimeYAML:  "name: my-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n",
			query:        "my-model",
			wantName:     "my-model-id",
		},
		{
			name:         "found by name",
			activeEngine: "my-engine",
			modelYAML:    "name: my-model-id\nalias: my-model\ndisk-size: 1G\n",
			engineYAML:   "name: my-engine\nruntime: my-runtime\nmodel:\n  options:\n    - my-model-id\n",
			runtimeYAML:  "name: my-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n",
			query:        "my-model-id",
			wantAlias:    "my-model",
		},
		{
			name:         "no active engine",
			activeEngine: "",
			query:        "my-model",
			wantErr:      true,
		},
		{
			name:         "incompatible with active engine",
			activeEngine: "my-engine",
			modelYAML:    "id: other-model-id\nname: other-model\ndisk-size: 1G\n",
			engineYAML:   "name: my-engine\nmodel:\n  options:\n    - some-other-model-id\n",
			runtimeYAML:  "name: my-runtime\n",
			query:        "other-model",
			wantErr:      true,
		},
		{
			name:         "model does not exist",
			activeEngine: "my-engine",
			engineYAML:   "name: my-engine\nmodel:\n  options: []\n",
			runtimeYAML:  "name: my-runtime\n",
			query:        "nonexistent-model",
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			modelsDir := t.TempDir()
			enginesDir := t.TempDir()
			runtimeDir := t.TempDir()

			if tc.modelYAML != "" {
				writeModelYAML(t, modelsDir, "my-model", tc.modelYAML)
			}
			if tc.engineYAML != "" {
				writeEngineYAML(t, enginesDir, "my-engine", tc.engineYAML)
			}
			if tc.runtimeYAML != "" {
				writeRuntimeYAML(t, runtimeDir, "my-runtime", tc.runtimeYAML)
			}
			cache := storage.NewMockCache()
			if err := cache.SetActiveEngine(tc.activeEngine); err != nil {
				t.Fatalf("SetActiveEngine: %v", err)
			}
			ctx := &Context{ModelsDir: modelsDir, EnginesDir: enginesDir, RuntimesDir: runtimeDir, Cache: cache}

			manifest, err := GetModelManifestByNameOrAlias(ctx, tc.query, getTestMachine())

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantName != "" && manifest.Name != tc.wantName {
				t.Errorf("expected Name %q, got %q", tc.wantName, manifest.Name)
			}
			if tc.wantAlias != "" && manifest.Alias != tc.wantAlias {
				t.Errorf("expected Alias %q, got %q", tc.wantAlias, manifest.Alias)
			}
		})
	}
}

func TestGetAllModels(t *testing.T) {
	modelsDir := t.TempDir()
	enginesDir := t.TempDir()
	runtimeDir := t.TempDir()
	writeModelYAML(t, modelsDir, "model1", "name: model1\nalias: m1\ndisk-size: 1G\n")
	writeModelYAML(t, modelsDir, "model2", "name: model2\nalias: m2\ndisk-size: 2G\n")
	writeRuntimeYAML(t, runtimeDir, "my-runtime", "name: my-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")
	writeEngineYAML(t, enginesDir, "my-engine", "name: my-engine\nruntime: my-runtime\nmodel:\n  options:\n    - model1\n")

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("my-engine"); err != nil {
		t.Fatalf("SetActiveEngine: %v", err)
	}
	ctx := &Context{ModelsDir: modelsDir, EnginesDir: enginesDir, RuntimesDir: runtimeDir, Cache: cache}

	models, err := GetAllModels(ctx, getTestMachine())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}

	expectedNames := map[string]bool{"model1": true, "model2": true}
	for _, model := range models {
		if !expectedNames[model.Name] {
			t.Errorf("unexpected model name: %s", model.Name)
		}
	}
}

func TestGetAllModelsWithEngines(t *testing.T) {
	modelsDir := t.TempDir()
	enginesDir := t.TempDir()
	runtimeDir := t.TempDir()

	writeModelYAML(t, modelsDir, "model1", "name: model1\nalias: m1\ndisk-size: 1G\n")
	writeModelYAML(t, modelsDir, "model2", "name: model2\nalias: m2\ndisk-size: 2G\n")
	writeModelYAML(t, modelsDir, "model3", "name: model3\nalias: m3\ndisk-size: 3G\n")
	writeModelYAML(t, modelsDir, "model4", "name: model4\nalias: m4\ndisk-size: 4G\n")
	writeRuntimeYAML(t, runtimeDir, "my-runtime", "name: my-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")

	writeEngineYAML(t, enginesDir, "engine1", "name: engine1\nruntime: my-runtime\nmodel:\n  options:\n    - model1\n    - model3\n")
	writeEngineYAML(t, enginesDir, "engine2", "name: engine2\nruntime: my-runtime\nmodel:\n  options:\n    - model1\n    - model2\n")

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("engine1"); err != nil {
		t.Fatalf("SetActiveEngine: %v", err)
	}
	ctx := &Context{ModelsDir: modelsDir, EnginesDir: enginesDir, RuntimesDir: runtimeDir, Cache: cache}

	modelsWithEngines, err := GetAllModels(ctx, getTestMachine())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modelsWithEngines) != 4 {
		t.Fatalf("expected 4 models, got %d", len(modelsWithEngines))
	}

	for _, modelWithEngines := range modelsWithEngines {
		switch modelWithEngines.Name {
		case "model1":
			if len(modelWithEngines.CompatibleEngines) != 2 || modelWithEngines.CompatibleEngines[0] != "engine1" || modelWithEngines.CompatibleEngines[1] != "engine2" {
				t.Errorf("model1 should be compatible with engine1 and engine2, got: %v", modelWithEngines.CompatibleEngines)
			}
		case "model2":
			if len(modelWithEngines.CompatibleEngines) != 1 || modelWithEngines.CompatibleEngines[0] != "engine2" {
				t.Errorf("model2 should be compatible with engine2, got: %v", modelWithEngines.CompatibleEngines)
			}
		case "model3":
			if len(modelWithEngines.CompatibleEngines) != 1 || modelWithEngines.CompatibleEngines[0] != "engine1" {
				t.Errorf("model3 should be compatible with engine1, got: %v", modelWithEngines.CompatibleEngines)
			}
		case "model4":
			if len(modelWithEngines.CompatibleEngines) != 0 {
				t.Errorf("model4 should not be compatible with any engines, got: %v", modelWithEngines.CompatibleEngines)
			}
		default:
			t.Errorf("unexpected model name: %s", modelWithEngines.Name)
		}
	}
}

func TestAvailableDiskSpace(t *testing.T) {
	machine := new(machine.Machine{})
	machine.Disk = []disk.Disk{
		{
			Available:  1024,
			Total:      2048,
			Path:       "/var/lib/snapd/snaps",
			MountPoint: new("/"),
		},
	}

	diskSpace, err := availableDiskSpace(machine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diskSpace != 1024 {
		t.Errorf("expected positive available disk space, got: %d", diskSpace)
	}
}

func TestAvailableMemory(t *testing.T) {
	machine := machine.Machine{
		CPUs: []cpu.CPU{{
			Architecture:   "amd64",
			ManufacturerId: "GenuineIntel",
			Flags:          []string{"fpu", "vme", "de"},
		}},
		Memory: memory.Memory{TotalRam: 67012501504, TotalSwap: 0}, // 64 GiB RAM, no swap
		Disk: []disk.Disk{{
			Total:     1006451294208, // ~937 GiB
			Available: 943543738368,  // ~878 GiB
			Path:      "/var/lib/snapd/snaps",
		}},
		PCIDevices: []pci.Device{
			{
				Bus:                  "pci",
				Slot:                 "0000:00:00.0",
				BusNumber:            0x0,
				DeviceClass:          0x380,
				ProgrammingInterface: new(uint8(0)),
				VendorId:             0x1002, //amd
				DeviceId:             0x4637,
				SubvendorId:          new(uint16(0x103C)),
				SubdeviceId:          new(uint16(0x89C6)),
				AdditionalProperties: map[string]string{
					"vram":              "10737418240", // 10 GiB
					"microarchitecture": "gfx1153",
				},
			},
			{
				Bus:                  "pci",
				Slot:                 "0000:00:00.0",
				BusNumber:            0x0,
				DeviceClass:          0x380,
				ProgrammingInterface: new(uint8(0)),
				VendorId:             0x1002, //amd
				DeviceId:             0x4637,
				SubvendorId:          new(uint16(0x103C)),
				SubdeviceId:          new(uint16(0x89C6)),
				AdditionalProperties: map[string]string{
					"vram":              "107374182400", // 100 GiB
					"microarchitecture": "gfx1152",
				},
			},
			{
				Bus:                  "pci",
				Slot:                 "0000:00:00.0",
				BusNumber:            0x0,
				DeviceClass:          0x380,
				ProgrammingInterface: new(uint8(0)),
				VendorId:             0x10de, // nvidia
				DeviceId:             0x4637,
				SubvendorId:          new(uint16(0x103C)),
				SubdeviceId:          new(uint16(0x89C6)),
				AdditionalProperties: map[string]string{
					"vram":               "[N/A]", // 10 GiB
					"compute-capability": "6.7",
				},
			},
		},
	}

	//test that the vram is computed for the correct device
	engine, err := engines.LoadManifest("../../../test_data/engines", "rocm-generic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	memoryAvailable, err := availableMemory(&machine, *engine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memoryAvailable != 107374182400 {
		t.Errorf("expected 107374182400 available memory, got: %d", memoryAvailable)
	}

	// test that vram is ignored if engine is not gpu capable
	engine, err = engines.LoadManifest("../../../test_data/engines", "cpu")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	memoryAvailable, err = availableMemory(&machine, *engine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memoryAvailable != 67012501504 {
		t.Errorf("expected 67012501504 available memory, got: %d", memoryAvailable)
	}

	// test that system memory is returned when nvidia-smi return "[N/A]" for the VRAM
	engine, err = engines.LoadManifest("../../../test_data/engines", "cuda-generic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	memoryAvailable, err = availableMemory(&machine, *engine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memoryAvailable != 67012501504 {
		t.Errorf("expected 67012501504 available memory, got: %d", memoryAvailable)
	}
}
