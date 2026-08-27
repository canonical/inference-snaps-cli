package fastrpc

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func TestMatch(t *testing.T) {
	hostDevices := []types.DetectedDevice{
		{
			Type: "npu",
			Bus:  "fastrpc",
			Metadata: &types.DeviceMetadata{
				VendorName: "qualcomm",
				ProductName:   "/dev/fastrpc-cdsp",
			},
		},
		{
			Type: "npu",
			Bus:  "fastrpc",
			Metadata: &types.DeviceMetadata{
				VendorName: "qualcomm",
				ProductName:   "/dev/fastrpc-cdsp1-secure",
			},
		},
	}

	t.Run("FastRPC bus", func(t *testing.T) {
		device := engines.Device{Type: "npu", Bus: "fastrpc"}
		score, issues := Match(device, hostDevices)
		if score == 0 || len(issues) > 0 {
			t.Fatalf("expected a positive score with no issues, score=%d issues=%v", score, issues)
		}
	})

	t.Run("custom node glob", func(t *testing.T) {
		nodeGlob := "/dev/fastrpc-cdsp1*"
		device := engines.Device{Type: "npu", Bus: "fastrpc", NodeGlob: &nodeGlob}
		score, issues := Match(device, hostDevices)
		if score == 0 || len(issues) > 0 {
			t.Fatalf("expected a positive score with no issues, score=%d issues=%v", score, issues)
		}
	})

	t.Run("no matching node glob", func(t *testing.T) {
		nodeGlob := "/dev/fastrpc-gdsp*"
		device := engines.Device{Type: "npu", Bus: "fastrpc", NodeGlob: &nodeGlob}
		score, issues := Match(device, hostDevices)
		if score != 0 {
			t.Fatalf("expected score=0, got %d", score)
		}
		if len(issues) == 0 {
			t.Fatal("expected compatibility issues")
		}
	})

	t.Run("FastRPC device without legacy metadata", func(t *testing.T) {
		device := engines.Device{Type: "npu", Bus: "fastrpc"}
		host := []types.DetectedDevice{{Bus: "fastrpc"}}
		score, issues := Match(device, host)
		if score == 0 || len(issues) > 0 {
			t.Fatalf("expected a positive score with no issues, score=%d issues=%v", score, issues)
		}
	})
}
