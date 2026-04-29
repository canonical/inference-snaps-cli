package renesas

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func TestInfo(t *testing.T) {
	tmp := t.TempDir()
	originalRoot := devRoot
	devRoot = tmp
	t.Cleanup(func() {
		devRoot = originalRoot
	})

	_, err := os.Create(filepath.Join(tmp, "drpai1"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmp, "drpai0"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmp, "something-else"))
	if err != nil {
		t.Fatal(err)
	}

	devices, err := Info()
	if err != nil {
		t.Fatal(err)
	}

	expectedDevices := []types.DetectedDevice{
		{
			Type: "npu",
			Bus:  "drpai",
			Metadata: &types.DeviceMetadata{
				VendorName:  "renesas",
				ProductName: filepath.Join(tmp, "drpai0"),
			},
		},
		{
			Type: "npu",
			Bus:  "drpai",
			Metadata: &types.DeviceMetadata{
				VendorName:  "renesas",
				ProductName: filepath.Join(tmp, "drpai1"),
			},
		},
	}
	if !reflect.DeepEqual(devices, expectedDevices) {
		t.Fatalf("devices=%+v expected=%+v", devices, expectedDevices)
	}
}

func TestDetectFromNodes(t *testing.T) {
	nodes := []string{
		"/dev/drpai1",
		"/dev/something-else",
		"/dev/drpai0",
	}

	devices := DetectFromNodes(nodes)

	expectedDevices := []types.DetectedDevice{
		{
			Type: "npu",
			Bus:  "drpai",
			Metadata: &types.DeviceMetadata{
				VendorName:  "renesas",
				ProductName: "/dev/drpai0",
			},
		},
		{
			Type: "npu",
			Bus:  "drpai",
			Metadata: &types.DeviceMetadata{
				VendorName:  "renesas",
				ProductName: "/dev/drpai1",
			},
		},
	}

	if !reflect.DeepEqual(devices, expectedDevices) {
		t.Fatalf("devices=%+v expected=%+v", devices, expectedDevices)
	}
}
