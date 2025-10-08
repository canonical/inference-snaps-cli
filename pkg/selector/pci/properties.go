package pci

import (
	"fmt"
	"strconv"

	"github.com/canonical/famous-models-cli/pkg/engines"
	"github.com/canonical/famous-models-cli/pkg/selector/weights"
	"github.com/canonical/famous-models-cli/pkg/types"
	"github.com/canonical/famous-models-cli/pkg/utils"
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

	// TODO compute_capability

	// has_driver - if this field is set, we can definitely say the driver is, or is not available
	if hasDriverStr, ok := pciDevice.AdditionalProperties["has_driver"]; ok {
		hasDriver, err := strconv.ParseBool(hasDriverStr)
		if err != nil {
			return 0, reasons, fmt.Errorf("error parsing has_driver property: %v", err)
		}
		if !hasDriver {
			reasons = append(reasons, "device driver not installed")
			return 0, reasons, nil
		} else {
			extraScore += weights.HasDriver
		}
	}

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
