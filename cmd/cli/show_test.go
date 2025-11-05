package main

import (
	"encoding/json"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"gopkg.in/yaml.v3"
)

func TestInfoLong(t *testing.T) {
	engine, err := selector.LoadManifestFromDir("../../test_data/engines", "intel-gpu")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	showEngineFormat = "yaml"
	err = printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoShort(t *testing.T) {
	engine, err := selector.LoadManifestFromDir("../../test_data/engines", "cpu-avx1")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	showEngineFormat = "yaml"
	err = printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoJSON(t *testing.T) {
	engine, err := selector.LoadManifestFromDir("../../test_data/engines", "intel-gpu")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	showEngineFormat = "json"
	err = printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the output is valid JSON by attempting to marshal/unmarshal
	jsonBytes, err := json.MarshalIndent(scoredEngine, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal engine to JSON: %v", err)
	}

	var unmarshaled engines.ScoredManifest
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
}

func TestInfoYAML(t *testing.T) {
	engine, err := selector.LoadManifestFromDir("../../test_data/engines", "cpu-avx1")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	showEngineFormat = "yaml"
	err = printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the output is valid YAML by attempting to marshal/unmarshal
	yamlBytes, err := yaml.Marshal(scoredEngine)
	if err != nil {
		t.Fatalf("failed to marshal engine to YAML: %v", err)
	}

	var unmarshaled engines.ScoredManifest
	err = yaml.Unmarshal(yamlBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}
}

func TestInfoInvalidFormat(t *testing.T) {
	engine, err := selector.LoadManifestFromDir("../../test_data/engines", "cpu-avx1")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	showEngineFormat = "invalid"
	err = printEngineManifest(scoredEngine)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}
