package mediatek

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func TestDetectFromCompatibles(t *testing.T) {
	platforms, devices := DetectFromCompatibles([]string{
		"arm,v8",
		"mediatek,mt8195-apusys_rv",
		"mediatek,mt8188-apusys_rv",
		"mediatek,mt8189-apusys_rv",
		"mediatek,mt9999-apusys_rv",
		"mediatek,mt8188-apusys_rv", // duplicate should be deduped
	})

	expectedPlatforms := []types.PlatformInfo{
		{Vendor: "mediatek", Name: "genio-1200"},
		{Vendor: "mediatek", Name: "genio-510-700"},
		{Vendor: "mediatek", Name: "genio-520-720"},
		{Vendor: "mediatek", Name: "mediatek-mt9999"},
	}
	if !reflect.DeepEqual(platforms, expectedPlatforms) {
		t.Fatalf("platforms=%+v expected=%+v", platforms, expectedPlatforms)
	}

	expectedDevices := []types.DetectedDevice{{
		Type: "npu",
		Bus:  "mdla",
		Nodes: []string{
			"mediatek,mt8188-apusys_rv",
			"mediatek,mt8189-apusys_rv",
			"mediatek,mt8195-apusys_rv",
			"mediatek,mt9999-apusys_rv",
		},
	}}
	if !reflect.DeepEqual(devices, expectedDevices) {
		t.Fatalf("devices=%+v expected=%+v", devices, expectedDevices)
	}
}

func TestInfo(t *testing.T) {
	tmp := t.TempDir()
	originalRoot := deviceTreeRoot
	deviceTreeRoot = tmp
	t.Cleanup(func() {
		deviceTreeRoot = originalRoot
	})

	goodDir := filepath.Join(tmp, "soc", "apu")
	if err := os.MkdirAll(goodDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Device-tree compatible entries are NUL-separated.
	if err := os.WriteFile(filepath.Join(goodDir, "compatible"), []byte("arm,v8\x00mediatek,mt8188-apusys_rv\x00"), 0o644); err != nil {
		t.Fatal(err)
	}

	otherDir := filepath.Join(tmp, "soc", "cpu")
	if err := os.MkdirAll(otherDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(otherDir, "compatible"), []byte("arm,cortex-a76\x00"), 0o644); err != nil {
		t.Fatal(err)
	}

	platforms, devices, err := Info()
	if err != nil {
		t.Fatal(err)
	}

	expectedPlatforms := []types.PlatformInfo{{Vendor: "mediatek", Name: "genio-510-700"}}
	if !reflect.DeepEqual(platforms, expectedPlatforms) {
		t.Fatalf("platforms=%+v expected=%+v", platforms, expectedPlatforms)
	}

	expectedDevices := []types.DetectedDevice{{
		Type:  "npu",
		Bus:   "mdla",
		Nodes: []string{"mediatek,mt8188-apusys_rv"},
	}}
	if !reflect.DeepEqual(devices, expectedDevices) {
		t.Fatalf("devices=%+v expected=%+v", devices, expectedDevices)
	}
}

func TestInfoNoMatch(t *testing.T) {
	tmp := t.TempDir()
	originalRoot := deviceTreeRoot
	deviceTreeRoot = tmp
	t.Cleanup(func() {
		deviceTreeRoot = originalRoot
	})

	if err := os.MkdirAll(filepath.Join(tmp, "soc"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "soc", "compatible"), []byte("arm,v8\x00"), 0o644); err != nil {
		t.Fatal(err)
	}

	platforms, devices, err := Info()
	if err != nil {
		t.Fatal(err)
	}
	if len(platforms) != 0 || len(devices) != 0 {
		t.Fatalf("expected no mediatek npu match, got platforms=%+v devices=%+v", platforms, devices)
	}
}