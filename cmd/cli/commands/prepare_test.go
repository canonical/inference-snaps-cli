package commands

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestPrepareLoad(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("SNAP_COMMON", tmpDir)

	cmd := &prepareCommand{
		noRestart: true,
		install:   true,
		Context: &common.Context{
			Config:     storage.NewMockConfig(),
			EnginesDir: "../../../test_data/engines",
			Cache:      storage.NewMockCache(),
			Snap:       snap.Mock(),
		},
	}
	useEngineCmd := useEngineCommand{Context: cmd.Context, noRestart: cmd.noRestart}
	scoredEngines, err := common.ScoreEnginesWithSpinner(cmd.Context)
	if err != nil {
		t.Fatalf("scoring engines: %v", err)
	}
	// Get rid of every non cpu engine from scoredengines
	var cpuEngines []engines.ScoredManifest
	for _, engine := range scoredEngines {
		if engine.Name == "cpu" || engine.Name == "cpu-devel" {
			cpuEngines = append(cpuEngines, engine)
		}
	}

	if err := cmd.load("../../../snap/configurations.yaml",
		cmd.Context, cpuEngines, &useEngineCmd); err != nil {
		t.Fatalf("unexpected error loading configurations: %v", err)
	}
}

func TestPrepareSetConfigurationsInvalidConfigType(t *testing.T) {
	cmd := &prepareCommand{
		Context: &common.Context{
			Config: storage.NewMockConfig(),
		},
	}

	changed, err := cmd.setConfigurations(map[string]any{
		"invalid": map[string]any{"key": "value"},
	})
	if err == nil {
		t.Fatal("expected an error for an unknown configuration type")
	}
	if changed {
		t.Fatal("did not expect changes when configuration type is invalid")
	}

	if got, want := err.Error(), `unknown configuration type "invalid"`; got != want {
		t.Fatalf("unexpected error: got %q, want %q", got, want)
	}
}

func TestPrepareSetConfigurationsChanged(t *testing.T) {
	cmd := &prepareCommand{
		Context: &common.Context{
			Config: storage.NewMockConfig(),
		},
	}

	changed, err := cmd.setConfigurations(map[string]any{
		"user": map[string]any{"foo": "bar"},
	})
	if err != nil {
		t.Fatalf("unexpected error setting configurations: %v", err)
	}
	if !changed {
		t.Fatal("expected changes when setting a new configuration value")
	}

	changed, err = cmd.setConfigurations(map[string]any{
		"user": map[string]any{"foo": "bar"},
	})
	if err != nil {
		t.Fatalf("unexpected error setting same configurations: %v", err)
	}
	if changed {
		t.Fatal("did not expect changes when re-applying the same configuration value")
	}
}
