package commands

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestSetValueValidation(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{})
	cmd := setCommand{
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}
	tests := map[string]struct {
		input       string
		errContains string
	}{
		"missing equal sign": {
			input:       "model",
			errContains: "expected key=value",
		},
		"starts with equal sign": {
			input:       "=value",
			errContains: "key must not start with an equal sign",
		},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			err := cmd.setValue(testCase.input)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", testCase.errContains)
			}
			if !strings.Contains(err.Error(), testCase.errContains) {
				t.Fatalf("expected error containing %q, got %q", testCase.errContains, err.Error())
			}
		})
	}
}

func TestSetValueSuccessForUserConfig(t *testing.T) {
	config := storage.NewMockConfig(map[string]any{})
	cmd := setCommand{
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	err := cmd.setValue("api.endpoint=https://example.com?x=1=y")
	if err != nil {
		t.Fatalf("setValue returned an unexpected error: %v", err)
	}

	values, err := config.Get("api.endpoint")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value, found := values["api.endpoint"]; !found || value != "https://example.com?x=1=y" {
		t.Fatalf("expected api.endpoint to be set to full value, got %#v", values)
	}
}
