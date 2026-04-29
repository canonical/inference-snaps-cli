package selector

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func TestDrpAiNpuSelection(t *testing.T) {
	hwInfo := &types.HwInfo{
		Memory: types.MemoryInfo{TotalRam: 8 * 1024 * 1024 * 1024},
		Disk: map[string]types.DirStats{
			"/var/lib/snapd/snaps": {
				Avail: 20 * 1024 * 1024 * 1024,
			},
		},
		Devices: []types.DetectedDevice{
			{
				Type: "npu",
				Bus:  "drpai",
				Metadata: &types.DeviceMetadata{
					VendorName:  "renesas",
					ProductName: "/dev/drpai0",
				},
			},
		},
	}

	manifest := engines.Manifest{
		Name:        "renesas-npu",
		Description: "renesas drp-ai npu",
		Vendor:      "renesas",
		Grade:       "stable",
		Devices: engines.Devices{
			Anyof: []engines.Device{{
				Type: "npu",
				Bus:  "drpai",
			}},
		},
	}

	score, report, err := checkEngine(hwInfo, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !report.EngineCompatible() {
		t.Fatalf("expected engine to be compatible: %+v", report)
	}
	if score <= 0 {
		t.Fatalf("expected positive score, got %d", score)
	}
}
