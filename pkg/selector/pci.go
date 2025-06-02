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

		vendorFound, err := checkPciVendorId(device, peripheral)
		if err != nil {
			return 0, err
		}
		if vendorFound {
			deviceScore += WeightPciVendor

			// Only check the device ID if the vendor is a match - device IDs are only unique per vendor namespace
			modelFound, err := checkPciDeviceId(device, peripheral)
			if err != nil {
				return 0, err
			}
			if modelFound {
				deviceScore += WeightPciModel
			}
		}
	}

	return deviceScore, nil
}

func checkPciVendorId(device types.StackDevice, peripheral types.PciPeripheral) (bool, error) {
	if device.VendorId != nil && peripheral.VendorId != "" {
		deviceVendorId, err := strconv.ParseUint(peripheral.VendorId, 0, 32)
		if err != nil {
			return false, fmt.Errorf("can't parse stack device vendor ID: %v", err)
		}
		peripheralVendorId, err := strconv.ParseUint(peripheral.VendorId, 0, 32)
		if err != nil {
			return false, fmt.Errorf("can't parse peripheral vendor ID: %v", err)
		}
		if deviceVendorId == peripheralVendorId {
			return true, nil
		}
	}
	return false, fmt.Errorf("pci device with required vendor not found")
}

func checkPciDeviceId(device types.StackDevice, peripheral types.PciPeripheral) (bool, error) {
	// For now, we use the model-ids list in the stack manifest to match the PCI device ID
	if len(device.ModelIds) > 0 && peripheral.DeviceId != "" {
		var deviceModelIds []uint64
		for _, modelIdString := range device.ModelIds {
			modelIdInt, err := strconv.ParseUint(modelIdString, 0, 32)
			if err != nil {
				return false, fmt.Errorf("can't parse stack device model ID: %v", err)
			}
			deviceModelIds = append(deviceModelIds, modelIdInt)
		}

		peripheralModelId, err := strconv.ParseUint(peripheral.DeviceId, 0, 32)
		if err != nil {
			return false, fmt.Errorf("can't parse peripheral device ID: %v", err)
		}

		if slices.Contains(deviceModelIds, peripheralModelId) {
			return true, nil
		}
	}
	return false, fmt.Errorf("pci device with required device id (model id) not found")
}
