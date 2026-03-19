package amd

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

var hwInfoGpu = types.PciDevice{
	Slot: "0000:c4:00.0",
}

func TestGfxArchitecture(t *testing.T) {
	t.Run("gfxArchitecture", func(t *testing.T) {
		gfxVersion, err := gfxArchitecture(hwInfoGpu, "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("gfx architecture: %v", gfxVersion)
	})
}

func TestGetGfxTargetVersion(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		errContains   string
		expectFailure bool
	}{
		{
			name:     "valid gfx target version",
			input:    "gfx_target_version 110002",
			expected: "gfx1102",
		},
		{
			name:          "invalid zero value",
			input:         "gfx_target_version 0",
			errContains:   "gfx_target_version is invalid for this device",
			expectFailure: true,
		},
		{
			name:          "unexpected format missing value",
			input:         "gfx_target_version",
			errContains:   "unexpected format for gfx_target_version",
			expectFailure: true,
		},
		{
			name:          "unexpected format non numeric",
			input:         "gfx_target_version abcdef",
			errContains:   "unexpected format for gfx_target_version",
			expectFailure: true,
		},
		{
			name:          "unexpected short numeric format",
			input:         "gfx_target_version 12345",
			errContains:   "gfx_target_version has an unexpected format",
			expectFailure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGfxTargetVersion(tt.input)
			t.Logf("input=%q got=%q err=%v", tt.input, got, err)

			if tt.expectFailure {
				if err == nil {
					t.Fatalf("expected error, got nil (result: %q)", got)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGpuProperties(t *testing.T) {
	originalFunc := utils.GetGlobalRootDir
	defer func() { utils.GetGlobalRootDir = originalFunc }() // Restore the original function after test

	mockCalled := false
	utils.GetGlobalRootDir = func() string {
		mockCalled = true
		return "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/"
	}

	t.Run("gpuProperties", func(t *testing.T) {
		properties, err := gpuProperties(hwInfoGpu)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("GPU properties: %v", properties)
	})
	if !mockCalled {
		t.Error("mock function was not called")
	}
}
