package selector

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func checkPciPeripherals(peripherals []types.PciPeripheral, device types.StackDevice) (int, error) {
	deviceScore := 0

	for _, peripheral := range peripherals {

		// Only check pci device vendor if it is set in the stack manifest
		if device.VendorId != nil {
			vendorMatch, err := checkPciVendorId(device, peripheral)
			if err != nil {
				return 0, fmt.Errorf("failed to check pci device vendor: %v", err)
			}
			if vendorMatch {
				deviceScore += WeightPciVendor

				// Only check the device ID if the vendor is a match - device IDs are only unique per vendor namespace
				if len(device.ModelIds) > 0 {
					modelMatch, err := checkPciDeviceId(device, peripheral)
					if err != nil {
						return 0, fmt.Errorf("failed to check pci device/model id: %v", err)
					}
					if modelMatch {
						deviceScore += WeightPciModel
					}
				}
			}
		}

		// TODO other PCI device fields
	}

	return deviceScore, nil
}

func checkPciVendorId(device types.StackDevice, peripheral types.PciPeripheral) (bool, error) {
	if peripheral.VendorId != "" {
		stackVendorId, err := strconv.ParseUint(*device.VendorId, 0, 32)
		if err != nil {
			return false, fmt.Errorf("can't parse stack device vendor ID: %v", err)
		}

		hwVendorId, err := strconv.ParseUint(peripheral.VendorId, 0, 32)
		if err != nil {
			return false, fmt.Errorf("can't parse peripheral vendor ID: %v", err)
		}

		if stackVendorId == hwVendorId {
			return true, nil
		}

		return false, nil // vendor does not match
	}
	return false, fmt.Errorf("pci device vendor id not set")
}

func checkPciDeviceId(device types.StackDevice, peripheral types.PciPeripheral) (bool, error) {
	// For now, we use the model-ids list in the stack manifest to match the PCI device ID
	if peripheral.DeviceId != "" {
		var stackModelIds []uint64
		for _, modelIdString := range device.ModelIds {
			modelIdInt, err := strconv.ParseUint(modelIdString, 0, 32)
			if err != nil {
				return false, fmt.Errorf("can't parse stack device model ID: %v", err)
			}
			stackModelIds = append(stackModelIds, modelIdInt)
		}

		hwModelId, err := strconv.ParseUint(peripheral.DeviceId, 0, 32)
		if err != nil {
			return false, fmt.Errorf("can't parse peripheral device ID: %v", err)
		}

		if slices.Contains(stackModelIds, hwModelId) {
			return true, nil
		}

		return false, nil // model does not match
	}
	return false, fmt.Errorf("pci device id not set")
}
