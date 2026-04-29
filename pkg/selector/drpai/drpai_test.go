package drpai

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func TestMatch(t *testing.T) {
	hostDevices := []types.DetectedDevice{
		{
			Type: "npu",
			Bus:  "drpai",
			Metadata: &types.DeviceMetadata{
				VendorName:  "renesas",
				ProductName: "/dev/drpai0",
			},
		},
	}

	device := engines.Device{Type: "npu", Bus: "drpai"}
	score, issues := Match(device, hostDevices)
	if len(issues) > 0 {
		t.Fatalf("expected match, got issues=%v", issues)
	}
	if score <= 0 {
		t.Fatalf("expected positive score, got %d", score)
	}
}

func TestMatchWithNodeGlob(t *testing.T) {
	hostDevices := []types.DetectedDevice{
		{
			Type: "npu",
			Bus:  "drpai",
			Metadata: &types.DeviceMetadata{
				VendorName:  "renesas",
				ProductName: "/dev/drpai0",
			},
		},
	}

	nodeGlob := "/dev/drpai1"
	device := engines.Device{Type: "npu", Bus: "drpai", NodeGlob: &nodeGlob}
	score, issues := Match(device, hostDevices)
	if score != 0 {
		t.Fatalf("expected zero score, got %d", score)
	}
	if len(issues) == 0 || !strings.Contains(issues[0], "device node matching") {
		t.Fatalf("expected node matching issue, got %v", issues)
	}
}

func TestMatchInvalidNodeGlob(t *testing.T) {
	hostDevices := []types.DetectedDevice{
		{
			Type: "npu",
			Bus:  "drpai",
			Metadata: &types.DeviceMetadata{
				VendorName:  "renesas",
				ProductName: "/dev/drpai0",
			},
		},
	}

	nodeGlob := "["
	device := engines.Device{Type: "npu", Bus: "drpai", NodeGlob: &nodeGlob}
	score, issues := Match(device, hostDevices)
	if score != 0 {
		t.Fatalf("expected zero score, got %d", score)
	}
	if len(issues) == 0 || !strings.Contains(issues[0], "invalid node-glob") {
		t.Fatalf("expected invalid node glob issue, got %v", issues)
	}
}
