package pci

import (
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

func hasAdditionalProperties(device engines.Device) bool {
	if device.VRam != nil {
		return true
	}
	if device.ComputeCapability != nil {
		return true
	}

	return false
}

func checkProperties(device engines.Device, pciDevice types.PciDevice) (int, []string, error) {
	var reasons []string
	extraScore := 0

	// vram
	if device.VRam != nil {
		vramScore, vramReasons, err := checkVram(device, pciDevice)
		reasons = append(reasons, vramReasons...)
		if err != nil {
			return 0, reasons, err
		}
		if vramScore > 0 {
			extraScore += vramScore
		} else {
			return 0, reasons, nil
		}
	}

	// TODO compute-capability

	return extraScore, reasons, nil
}

func checkVram(device engines.Device, pciDevice types.PciDevice) (int, []string, error) {
	var reasons []string

	vramRequired, err := utils.StringToBytes(*device.VRam)
	if err != nil {
		return 0, reasons, err
	}
	if vram, ok := pciDevice.AdditionalProperties["vram"]; ok {
		vramAvailable, err := utils.StringToBytes(vram)
		if err != nil {
			return 0, reasons, err
		}
		if vramAvailable >= vramRequired {
			return weights.GpuVRam, reasons, nil
		} else {
			reasons = append(reasons, "not enough vram")
			return 0, reasons, nil
		}
	} else {
		// Hardware Info does not list available vram
		reasons = append(reasons, "hw-info missing additional properties field \"vram\"")
		return 0, reasons, nil
	}
}
