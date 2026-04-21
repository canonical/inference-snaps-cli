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

	platforms, devices, err := Info()
	if err != nil {
		t.Fatal(err)
	}

	expectedPlatforms := []types.PlatformInfo{{
		Vendor: "qualcomm",
		Name:   "dragonwing",
	}}
	if !reflect.DeepEqual(platforms, expectedPlatforms) {
		t.Fatalf("platforms=%+v expected=%+v", platforms, expectedPlatforms)
	}

	expectedDevices := []types.DetectedDevice{{
		Type: "npu",
		Bus:  "fastrpc",
		Nodes: []string{
			filepath.Join(tmp, "fastrpc-cdsp"),
			filepath.Join(tmp, "fastrpc-cdsp1"),
		},
	}}
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

	platforms, devices := DetectFromNodes(nodes)
	if len(platforms) != 1 {
		t.Fatalf("expected 1 platform, got %d", len(platforms))
	}
	if platforms[0].Vendor != "qualcomm" || platforms[0].Name != "dragonwing" {
		t.Fatalf("unexpected platform %+v", platforms[0])
	}

	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	if devices[0].Type != "npu" || devices[0].Bus != "fastrpc" {
		t.Fatalf("unexpected device %+v", devices[0])
	}

	expectedNodes := []string{"/dev/fastrpc-cdsp", "/dev/fastrpc-cdsp1-secure"}
	if !reflect.DeepEqual(devices[0].Nodes, expectedNodes) {
		t.Fatalf("nodes=%v expected=%v", devices[0].Nodes, expectedNodes)
	}
}
