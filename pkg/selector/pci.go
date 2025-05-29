package selector

import "github.com/canonical/ml-snap-utils/pkg/types"

func checkPciPeripherals(peripherals []types.PciPeripheral, device types.StackDevice) (int, error) {
	deviceScore := 0
	for _, peripheral := range peripherals {
		/*
			DeviceClass          string  `json:"device_class"`
			ProgrammingInterface *string `json:"programming_interface"`
			VendorId             string  `json:"vendor_id"`
			DeviceId             string  `json:"device_id"`
			SubvendorId          *string `json:"subvendor_id,omitempty"`
			SubdeviceId          *string `json:"subdevice_id,omitempty"`
		*/
		if device.VendorId != nil && peripheral.VendorId != "" {
			if *device.VendorId == peripheral.VendorId {
				deviceScore += WeightPciVendor
			}
		}
	}

	return deviceScore, nil
}
