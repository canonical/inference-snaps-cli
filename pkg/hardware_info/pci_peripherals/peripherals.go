package pci_peripherals

import (
	"fmt"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Info(friendlyNames bool) ([]types.PciPeripheral, error) {
	var peripherals []types.PciPeripheral

	pciDevices, err := pci.PciDevices(friendlyNames)
	if err != nil {
		return nil, err
	}

	// Filter the list of PCI devices and only keep peripherals that are of interest
	pciAccelerators, err := accelerators(pciDevices)
	if err != nil {
		return nil, err
	}
	peripherals = append(peripherals, pciAccelerators...)

	pciCoprocessors, err := coprocessors(pciDevices)
	if err != nil {
		return nil, err
	}
	peripherals = append(peripherals, pciCoprocessors...)

	return peripherals, err
}

func pciDeviceToPeripheral(device pci.PciDevice) types.PciPeripheral {
	peripheral := types.PciPeripheral{}

	peripheral.DeviceClass = fmt.Sprintf("0x%04x", device.DeviceClass)
	if device.ProgrammingInterface != nil {
		programmingInterface := fmt.Sprintf("0x%02x", *device.ProgrammingInterface)
		peripheral.ProgrammingInterface = &programmingInterface
	}

	peripheral.VendorId = fmt.Sprintf("0x%04x", device.VendorId)
	peripheral.DeviceId = fmt.Sprintf("0x%04x", device.DeviceId)
	if device.SubvendorId != nil {
		subVendorId := fmt.Sprintf("0x%04x", *device.SubvendorId)
		peripheral.SubvendorId = &subVendorId
	}
	if device.SubdeviceId != nil {
		subDeviceId := fmt.Sprintf("0x%04x", *device.SubdeviceId)
		peripheral.SubdeviceId = &subDeviceId
	}

	peripheral.VendorName = device.VendorName
	peripheral.DeviceName = device.DeviceName
	peripheral.SubvendorName = device.SubvendorName
	peripheral.SubdeviceName = device.SubdeviceName

	return peripheral
}

/*
accelerators takes a slice of PCI devices, filters it and returns a slice of PCI accelerator devices
Class:	1200 # 0x12 = Processing Accelerator - Intel Lunar Lake NPU identifies as this class
*/
func accelerators(pciDevices []pci.PciDevice) ([]types.PciPeripheral, error) {
	var peripherals []types.PciPeripheral

	for _, device := range pciDevices {
		if device.DeviceClass&0xFF00 == 0x1200 {
			peripherals = append(peripherals, pciDeviceToPeripheral(device))
		}
	}
	return peripherals, nil

}

/*
coprocessors takes a slice of PCI devices, filters it and returns a slice of PCI coprocessors
Class: 0b40 # 0x0B = Processor, 0x40 = Co-Processor - Hailo PCI devices identify as this class
*/
func coprocessors(pciDevices []pci.PciDevice) ([]types.PciPeripheral, error) {
	var peripherals []types.PciPeripheral

	for _, device := range pciDevices {
		if device.DeviceClass == 0x0B40 {
			peripherals = append(peripherals, pciDeviceToPeripheral(device))
		}
	}
	return peripherals, nil
}
