package main

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
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
