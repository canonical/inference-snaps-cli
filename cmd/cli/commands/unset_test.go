package commands

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestUnsetValueRemovesUserConfigWithoutRestart(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{"api.endpoint": "https://example.com"})
	cmd := unsetCommand{
		noRestart: false,
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	if err := cmd.unsetValue("api.endpoint"); err != nil {
		t.Fatalf("unsetValue returned an unexpected error: %v", err)
	}

	values, err := config.Get("api.endpoint")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("expected api.endpoint to be removed, got %#v", values)
	}
}
