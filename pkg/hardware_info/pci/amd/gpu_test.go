package amd

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

var hwInfoGpu = types.PciDevice{
	Slot: "0000:c4:00.0",
}

func TestGfxArchitecture(t *testing.T) {
	t.Run("gfxArchitecture", func(t *testing.T) {
		gfxVersion, err := gfxArchitecture(hwInfoGpu)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("gfx architecture: %v", gfxVersion)
	})
}

func TestVRam(t *testing.T) {
	t.Run("vRam", func(t *testing.T) {
		vram, err := vRam(hwInfoGpu)
		if err != nil {
			t.Fatal(err)
		}
		if vram == nil {
			t.Fatal("vRAM is nil")
		}
		t.Logf("vRAM: %d", *vram)
	})
}

func TestGpuProperties(t *testing.T) {
	t.Run("gpuProperties", func(t *testing.T) {
		properties, err := gpuProperties(hwInfoGpu)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("GPU properties: %v", properties)
	})
}
