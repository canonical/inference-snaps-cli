package pci

import (
	"fmt"
	"testing"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func Match(device engines.Device, pcis []types.PciDevice) (maxDeviceScore int, deviceIssues []string) {
	maxDeviceScore = 0

	if len(pcis) == 0 {
		deviceIssues = append(deviceIssues, "no pci devices on host system")
		return
	}

	availableDevices := filterPciDevices(device, pcis)
	scoredDevices, scoreIssues := scorePciDevices(device, availableDevices)

	for _, pci := range scoredDevices {
		if pci.Score > maxDeviceScore {
			maxDeviceScore = pci.Score
		}
	}
	if maxDeviceScore == 0 {
		deviceIssues = append(deviceIssues, scoreIssues...)
		return
	}

	return
}

func filterPciDevices(device engines.Device, pcis []types.PciDevice) []types.PciDevice {
	var foundDevices []types.PciDevice
	for _, pciDevice := range pcis {
		include := true

		if device.VendorId != nil {
			if *device.VendorId != pciDevice.VendorId {
				include = false
			} else {
				// A model ID is only unique per vendor ID namespace. Only check it if the vendor is a match
				if device.DeviceId != nil {
					if *device.DeviceId != pciDevice.DeviceId {
						include = false
					}
				}
			}
		}

		if include {
			foundDevices = append(foundDevices, pciDevice)
		}
	}
	return foundDevices
}

func scorePciDevices(device engines.Device, pciDevices []types.PciDevice) ([]types.PciDevice, []string) {
	var issues []string

	if len(pciDevices) == 0 {
		issues = append(issues, "device not found")
	}

	for i, pciDevice := range pciDevices {
		deviceScore, deviceIssues := scorePciDevice(device, pciDevice)

		pciDevices[i].Score = deviceScore
		for _, issue := range deviceIssues {
			issues = append(issues, fmt.Sprintf("pci %s: %s", pciDevice.Slot, issue))
		}
	}
	return pciDevices, issues
}

func scorePciDevice(device engines.Device, pciDevice types.PciDevice) (deviceScore int, issues []string) {
	deviceScore = 0

	// Device type: tpu, npu, gpu, etc
	if device.Type != "" {
		match := checkType(device.Type, pciDevice)
		if match {
			deviceScore += weights.PciDeviceType
		} else {
			deviceScore = 0
			issues = append(issues, fmt.Sprintf("wrong device class 0x%04x", pciDevice.DeviceClass))
			return
		}
	}

	// Prefer dGPU above iGPU
	// PCI devices on bus 0 are considered internal, and anything else external/discrete
	if pciDevice.BusNumber > 0 {
		deviceScore += weights.PciDeviceExternal
	}

	// Check additional properties
	if hasAdditionalProperties(device) {
		propsScore, err := checkProperties(device, pciDevice)
		if err != nil {
			deviceScore = 0
			issues = append(issues, err.Error())
			return
		}
		deviceScore += propsScore
	}

	// Check drivers
	for _, connection := range device.SnapConnections {
		connected, err := checkSnapConnection(connection)
		if err != nil {
			deviceScore = 0
			issues = append(issues, fmt.Sprintf("error checking snap connection %q: %v", connection, err))
			return
		}
		if !connected {
			deviceScore = 0
			issues = append(issues, fmt.Sprintf("%q is not connected", connection))
			return
		}
	}

	return deviceScore, nil
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

func checkSnapConnection(connection string) (bool, error) {
	if testing.Testing() {
		// Tests do not necessarily run inside a snap
		// Stub out and always return true for all connections
		return true, nil
	}
	return snapctl.IsConnected(connection).Run()
}
