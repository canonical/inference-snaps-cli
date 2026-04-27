package qualcomm

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

	_, err := os.Create(filepath.Join(tmp, "fastrpc-adsp-secure"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmp, "fastrpc-cdsp"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmp, "fastrpc-cdsp1"))
	if err != nil {
		t.Fatal(err)
	}

	devices, err := Info()
	if err != nil {
		t.Fatal(err)
	}

	expectedDevices := []types.DetectedDevice{
		{
			Type: "NPU - adsp-secure",
			Bus:  "fastrpc",
			Metadata: &types.DeviceMetadata{
				Vendor: "qualcomm",
				Name:   filepath.Join(tmp, "fastrpc-adsp-secure"),
				VendorName:    "dragonwing",
			},
		},
		{
			Type: "NPU - cdsp",
			Bus:  "fastrpc",
			Metadata: &types.DeviceMetadata{
				Vendor: "qualcomm",
				Name:   filepath.Join(tmp, "fastrpc-cdsp"),
				VendorName:    "dragonwing",
			},
		},
		{
			Type: "NPU - cdsp1",
			Bus:  "fastrpc",
			Metadata: &types.DeviceMetadata{
				Vendor: "qualcomm",
				Name:   filepath.Join(tmp, "fastrpc-cdsp1"),
				VendorName:    "dragonwing",
			},
		},
	}
	if !reflect.DeepEqual(devices, expectedDevices) {
		t.Fatalf("devices=%+v expected=%+v", devices, expectedDevices)
	}
}

func TestDetectFromNodes(t *testing.T) {
	nodes := []string{
		"/dev/fastrpc-cdsp1-secure",
		"/dev/fastrpc-adsp",
		"/dev/something-else",
		"/dev/fastrpc-cdsp",
	}

	devices := DetectFromNodes(nodes)

	expectedDevices := []types.DetectedDevice{
		{
			Type: "NPU - adsp",
			Bus:  "fastrpc",
			Metadata: &types.DeviceMetadata{
				Vendor: "qualcomm",
				Name:   "/dev/fastrpc-adsp",
				VendorName:    "dragonwing",
			},
		},
		{
			Type: "NPU - cdsp",
			Bus:  "fastrpc",
			Metadata: &types.DeviceMetadata{
				Vendor: "qualcomm",
				Name:   "/dev/fastrpc-cdsp",
				VendorName:    "dragonwing",
			},
		},
		{
			Type: "NPU - cdsp1-secure",
			Bus:  "fastrpc",
			Metadata: &types.DeviceMetadata{
				Vendor: "qualcomm",
				Name:   "/dev/fastrpc-cdsp1-secure",
				VendorName:    "dragonwing",
			},
		},
	}

	if !reflect.DeepEqual(devices, expectedDevices) {
		t.Fatalf("devices=%+v expected=%+v", devices, expectedDevices)
	}
}
