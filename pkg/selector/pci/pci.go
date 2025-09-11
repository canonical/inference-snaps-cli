package pci

import (
	"fmt"

	"github.com/canonical/stack-utils/pkg/engines"
	"github.com/canonical/stack-utils/pkg/selector/weights"
	"github.com/canonical/stack-utils/pkg/types"
)

func Match(stackDevice engines.Device, pcis []types.PciDevice) (int, []string, error) {
	var reasons []string
	maxDeviceScore := 0

	for _, pciDevice := range pcis {
		deviceScore, deviceReasons, err := checkPciDevice(stackDevice, pciDevice)
		reasons = append(reasons, deviceReasons...)
		if err != nil {
			return 0, reasons, err
		}

		if deviceScore > 0 {
			if deviceScore > maxDeviceScore {
				maxDeviceScore = deviceScore
			}
		}
	}

	return maxDeviceScore, reasons, nil
}

func checkPciDevice(stackDevice engines.Device, pciDevice types.PciDevice) (int, []string, error) {
	var reasons []string
	currentDeviceScore := 0

	// Device type: tpu, npu, gpu, etc
	if stackDevice.Type != "" {
		match := checkType(stackDevice.Type, pciDevice)
		if match {
			currentDeviceScore += weights.PciDeviceType
		} else {
			reasons = append(reasons, fmt.Sprintf("pci device type mismatch: %s", stackDevice.Type))
			return 0, reasons, nil
		}
	}

	// Prefer dGPU above iGPU
	// PCI devices on bus 0 are considered internal, and anything else external/discrete
	if pciDevice.BusNumber > 0 {
		currentDeviceScore += weights.PciDeviceExternal
	}

	if stackDevice.VendorId != nil {
		if *stackDevice.VendorId == pciDevice.VendorId {
			currentDeviceScore += weights.PciVendorId
		} else {
			reasons = append(reasons, fmt.Sprintf("pci vendor id mismatch: %04x", *stackDevice.VendorId))
			return 0, reasons, nil
		}

		// A model ID is only unique per vendor ID namespace. Only check it if the vendor is a match
		if stackDevice.DeviceId != nil {
			if *stackDevice.DeviceId == pciDevice.DeviceId {
				currentDeviceScore += weights.PciDeviceId
			} else {
				reasons = append(reasons, fmt.Sprintf("pci device id mismatch: %04x", *stackDevice.DeviceId))
				return 0, reasons, nil
			}
		}
	}

	// Check additional properties
	if hasAdditionalProperties(stackDevice) {
		propsScore, propReasons, err := checkProperties(stackDevice, pciDevice)
		reasons = append(reasons, propReasons...)
		if err != nil {
			return 0, reasons, err
		}
		if propsScore > 0 {
			currentDeviceScore += propsScore
		} else {
			return 0, reasons, nil
		}
	}

	return currentDeviceScore, reasons, nil
}

func checkType(requiredType string, pciDevice types.PciDevice) bool {
	if requiredType == "gpu" {
		// 00 01 - legacy VGA devices
		// 03 xx - display controllers
		if pciDevice.DeviceClass == 0x0001 || pciDevice.DeviceClass&0xFF00 == 0x0300 {
			return true
		}
	}

	/*
		Base class 0x12 = Processing Accelerator - Intel Lunar Lake NPU identifies as this class
		Base class 0x0B = Processor, Sub class 0x40 = Co-Processor - Hailo PCI devices identify as this class
	*/
	if requiredType == "npu" || requiredType == "tpu" {
		if pciDevice.DeviceClass&0xFF00 == 0x1200 {
			// Processing accelerator
			return true
		}
		if pciDevice.DeviceClass == 0x0B40 {
			// Coprocessor
			return true
		}
	}

	return false
}
