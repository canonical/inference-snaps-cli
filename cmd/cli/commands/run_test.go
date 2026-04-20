package commands

import (
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestExtractPassthroughConfigs(t *testing.T) {
	configs := map[string]any{
		"passthrough.environment.my-key": "hello",
		"passthrough.other":              42,
		"regular.config":                 "ignored",
	}
	cmd := runCommand{}
	got, err := cmd.extractPassthroughConfigs(configs)
	if err != nil {
		t.Fatalf("extractPassthroughConfigs returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("extractPassthroughConfigs returned %d keys, want 2", len(got))
	}

	if got["environment.my-key"] != "hello" {
		t.Fatalf("extractPassthroughConfigs returned %v for environment.my-key, want hello", got["environment.my-key"])
	}

	if got["other"] != 42 {
		t.Fatalf("extractPassthroughConfigs returned %v for other, want 42", got["other"])
	}
}

func TestGetEnvVarsFromPassthroughConfigs(t *testing.T) {
	cmd := runCommand{}
	passthrough := map[string]any{
		"environment.my-key": "value",
		"environment.other":  123,
		"not-environment":    "ignored",
	}

	got, err := cmd.getEnvVarsFromPassthroughConfigs(passthrough)
	if err != nil {
		t.Fatalf("getEnvVarsFromPassthroughConfigs returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("getEnvVarsFromPassthroughConfigs returned %d keys, want 2", len(got))
	}

	if got["MY_KEY"] != "value" {
		t.Fatalf("getEnvVarsFromPassthroughConfigs returned %v for MY_KEY, want value", got["MY_KEY"])
	}

	if got["OTHER"] != "123" {
		t.Fatalf("getEnvVarsFromPassthroughConfigs returned %v for OTHER, want 123", got["OTHER"])
	}
}

func TestProcessPassthroughConfigs(t *testing.T) {
	cmd := runCommand{
		Context: &common.Context{
			Config: storage.NewMockConfig(map[string]any{
				"passthrough.environment.my-key": "value",
				"passthrough.environment.other":  123,
				"passthrough.not-environment":    "ignored",
			}),
		},
	}

	err := cmd.processPassthroughConfigs()
	if err != nil {
		t.Fatalf("processPassthroughConfigs returned error: %v", err)
	}

	if got := os.Getenv("MY_KEY"); got != "value" {
		t.Fatalf("expected MY_KEY to be %q, got %q", "value", got)
	}

	if got := os.Getenv("OTHER"); got != "123" {
		t.Fatalf("expected OTHER to be %q, got %q", "123", got)
	}
}
