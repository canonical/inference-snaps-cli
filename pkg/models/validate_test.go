package models

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func templateManifest() Manifest {
	manifest := Manifest{
		ID:           "test",
		Name:         "test",
		Description:  "test",
		ModelCardUrl: "https://example.com/model-card",
		Quantization: "Q4_K_M",
		Capabilities: []string{"text"},
		DiskSize:     "6G",
		Components:   []string{"test-component"},
		Environment:  []string{"MODEL_FILE=/tmp/model.gguf"},
	}
	return manifest
}

func TestManifestFiles(t *testing.T) {
	modelsDir := "../../test_data/models"

	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		t.Fatalf("Failed reading models directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			model := entry.Name()
			manifestPath := filepath.Join(modelsDir, model, ManifestFilename)
			t.Run(model, func(t *testing.T) {
				err = Validate(manifestPath)
				if err != nil {
					t.Fatalf("%s: %v", model, err)
				}
			})
		}
	}
}

func TestManifestEmpty(t *testing.T) {
	data := ""
	err := validateManifestYaml("", []byte(data))
	if err == nil {
		t.Fatal("Empty yaml should fail")
	}
	t.Log(err)
}

func TestUnknownField(t *testing.T) {
	data, _ := yaml.Marshal(templateManifest())
	data = append(data, []byte("unknown-field: test\n")...)

	err := validateManifestYaml("test", data)
	if err == nil {
		t.Fatal("Unknown field should fail")
	}
	t.Log(err)
}

func TestNameRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Name = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("name field is required")
	}
	t.Log(err)
}

func TestDescriptionRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Description = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("description field is required")
	}
	t.Log(err)
}

func TestModelCardUrlRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.ModelCardUrl = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("model-card-url field is required")
	}
	t.Log(err)
}

func TestQuantizationRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Quantization = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("quantization field is required")
	}
	t.Log(err)
}

func TestCapabilitiesRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Capabilities = nil

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("capabilities field is required")
	}
	t.Log(err)
}

func TestDiskSizeRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.DiskSize = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("disk-size field is required")
	}
	t.Log(err)
}

func TestComponentsRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Components = nil

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("components field is required")
	}
	t.Log(err)
}

func TestEnvironmentRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Environment = nil

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("environment field is required")
	}
	t.Log(err)
}

func TestModelIdMatch(t *testing.T) {
	manifest := templateManifest()
	manifest.ID = "different-id"

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("model directory name should match id in manifest")
	}
	t.Log(err)
}
