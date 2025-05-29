package pci_peripherals

import (
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Info(friendlyNames bool) ([]types.PciDevice, error) {
	var pciPeripherals []types.PciDevice

	pciDevices, err := pci.PciDevices(friendlyNames)
	if err != nil {
		return nil, err
	}

	// Filter the list of PCI devices and only keep peripherals that are of interest
	pciAccelerators, err := accelerators(pciDevices)
	if err != nil {
		return nil, err
	}
	pciPeripherals = append(pciPeripherals, pciAccelerators...)

	pciCoprocessors, err := coprocessors(pciDevices)
	if err != nil {
		return nil, err
	}
	pciPeripherals = append(pciPeripherals, pciCoprocessors...)

	return pciPeripherals, err
}

/*
accelerators takes a slice of PCI devices, filters it and returns a slice of PCI accelerator devices
Class:	1200 # 0x12 = Processing Accelerator - Intel Lunar Lake NPU identifies as this class
*/
func accelerators(pciDevices []types.PciDevice) ([]types.PciDevice, error) {
	var devices []types.PciDevice

	for _, device := range pciDevices {
		if device.DeviceClass&0xFF00 == 0x1200 {
			devices = append(devices, device)
		}
	}
	return devices, nil

}

/*
coprocessors takes a slice of PCI devices, filters it and returns a slice of PCI coprocessors
Class: 0b40 # 0x0B = Processor, 0x40 = Co-Processor - Hailo PCI devices identify as this class
*/
func coprocessors(pciDevices []types.PciDevice) ([]types.PciDevice, error) {
	var devices []types.PciDevice

	for _, device := range pciDevices {
		if device.DeviceClass == 0x0B40 {
			devices = append(devices, device)
		}
	}
	return devices, nil
}
