package commands

import (
	"os"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestParseKeyValue(t *testing.T) {
	cmd := setCommand{}

	tests := map[string]struct {
		input       string
		wantKey     string
		wantValue   string
		errContains string
	}{
		"empty input": {
			input:       "",
			errContains: "expected key=value",
		},
		"missing equal sign": {
			input:       "model",
			errContains: "expected key=value",
		},
		"starts with equal sign": {
			input:       "=value",
			errContains: "key must not start with an equal sign",
		},
		"simple pair": {
			input:     "model=llama",
			wantKey:   "model",
			wantValue: "llama",
		},
		"value keeps equal signs": {
			input:     "api.endpoint=https://example.com?a=b",
			wantKey:   "api.endpoint",
			wantValue: "https://example.com?a=b",
		},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			gotKey, gotValue, err := cmd.parseKeyValue(testCase.input)
			if testCase.errContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", testCase.errContains)
				}
				if !strings.Contains(err.Error(), testCase.errContains) {
					t.Fatalf("expected error containing %q, got %q", testCase.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("parseKeyValue returned an unexpected error: %v", err)
			}

			if gotKey != testCase.wantKey || gotValue != testCase.wantValue {
				t.Fatalf("expected (%q, %q), got (%q, %q)", testCase.wantKey, testCase.wantValue, gotKey, gotValue)
			}
		})
	}
}

func TestSetValueSuccessForUserConfig(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{"api.endpoint": "https://old.example.com"})
	cmd := setCommand{
		noRestart: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValues([]string{"api.endpoint=https://new.example.com"})
	if err != nil {
		t.Fatalf("setValue returned an unexpected error: %v", err)
	}

	values, err := config.Get("api.endpoint")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value, found := values["api.endpoint"]; !found || value != "https://new.example.com" {
		t.Fatalf("expected api.endpoint to be set to full value, got %#v", values)
	}
}

func TestSetValueRejectsUnknownKeys(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{})
	cmd := setCommand{
		noRestart: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValues([]string{"api.endpoint=https://example.com"})
	if err == nil {
		t.Fatal("expected error for unknown key, got nil")
	} else {
		if !strings.Contains(err.Error(), "unknown key") {
			t.Fatalf("expected unknown key error, got: %s", err)
		}
	}
}

func TestSetValuesSuccessForUserConfig(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{
		"api.endpoint": "https://old.example.com",
		"api.port":     8080,
	})
	cmd := setCommand{
		noRestart: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValues([]string{"api.endpoint=https://new.example.com", "api.port=9090"})
	if err != nil {
		t.Fatalf("setValues returned an unexpected error: %v", err)
	}

	values, err := config.Get("api")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value, found := values["api.endpoint"]; !found || value != "https://new.example.com" {
		t.Fatalf("expected api.endpoint to be updated, got %#v", values)
	}

	if value, found := values["api.port"]; !found || value != "9090" {
		t.Fatalf("expected api.port to be updated, got %#v", values)
	}
}

func TestSetValuesRejectsUnknownKeysAtomically(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{
		"api.endpoint": "https://old.example.com",
		"api.port":     8080,
	})
	cmd := setCommand{
		noRestart: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValues([]string{"api.endpoint=https://new.example.com", "unknown.key=value"})
	if err == nil {
		t.Fatal("expected unknown key error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown key") {
		t.Fatalf("expected unknown key error, got: %s", err)
	}

	values, err := config.Get("api")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value, found := values["api.endpoint"]; !found || value != "https://old.example.com" {
		t.Fatalf("expected no writes after validation error, got %#v", values)
	}
}

func TestSetValuesRejectsDuplicateKeys(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{"api.endpoint": "https://old.example.com"})
	cmd := setCommand{
		noRestart: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValues([]string{"api.endpoint=https://new.example.com", "api.endpoint=https://another.example.com"})
	if err == nil {
		t.Fatal("expected duplicate key error, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate key") {
		t.Fatalf("expected duplicate key error, got: %s", err)
	}

	values, err := config.Get("api.endpoint")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value, found := values["api.endpoint"]; !found || value != "https://old.example.com" {
		t.Fatalf("expected no writes after duplicate key error, got %#v", values)
	}
}

func TestSetAcceptsUnknownPassthroughKeys(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{})
	cmd := setCommand{
		noRestart: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValues([]string{"passthrough.custom-key=custom-value"})
	if err != nil {
		t.Fatalf("setValues returned an unexpected error for passthrough key: %v", err)
	}

	values, err := config.Get("passthrough.custom-key")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value, found := values["passthrough.custom-key"]; !found || value != "custom-value" {
		t.Fatalf("expected passthrough.custom-key to be set to custom-value, got %#v", values)
	}
}

func ExampleSet_assumeYesRestartServices() {
	if err := os.Setenv("SNAP_INSTANCE_NAME", "example-snap"); err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Unsetenv("SNAP_INSTANCE_NAME")
	}()

	config := storage.NewMockConfig(map[string]any{"api.endpoint": "https://old.example.com"})
	cmd := setCommand{
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	if err := cmd.setValues([]string{"api.endpoint=https://example.com"}); err != nil {
		panic(err)
	}

	// Output:
	// [mock] Restarting all services
}

func ExampleSet_noRestartWhenFinalValueUnchanged() {
	config := storage.NewMockConfig(map[string]any{"api.port": 8080})
	cmd := setCommand{
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValues([]string{"api.port=8080"})
	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleSet_restartWhenFinalValueUnchanged() {
	config := storage.NewMockConfig(map[string]any{"api.port": 8080})
	cmd := setCommand{
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValues([]string{"api.port=9999"})
	if err != nil {
		panic(err)
	}

	// Output:
	// [mock] Restarting all services
}
