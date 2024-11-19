package gpu

import (
	"fmt"

	"github.com/canonical/hardware-info/lspci"
)

func Info(friendlyNames bool) ([]Gpu, error) {
	pciDevices, err := lspci.PciDevices(friendlyNames)
	if err != nil {
		return nil, err
	}

	gpus, err := DisplayDevices(pciDevices)

	return gpus, nil
}

func DisplayDevices(pciDevices []lspci.PciDevice) ([]Gpu, error) {
	var gpus []Gpu

	for _, device := range pciDevices {
		// 00 01 - legacy VGA devices
		// 03 xx - display controllers
		if device.DeviceClass == 0x0001 || device.DeviceClass&0xFF00 == 0x0300 {
			var gpu Gpu
			gpu.VendorId = fmt.Sprintf("0x%04x", device.VendorId)
			gpu.DeviceId = fmt.Sprintf("0x%04x", device.DeviceId)
			if device.SubVendorId != nil {
				subVendorId := fmt.Sprintf("0x%04x", *device.SubVendorId)
				gpu.SubVendorId = &subVendorId
			}
			if device.SubDeviceId != nil {
				subDeviceId := fmt.Sprintf("0x%04x", *device.SubDeviceId)
				gpu.SubDeviceId = &subDeviceId
			}

			// TODO look up vram

			gpu.VendorName = device.VendorName
			gpu.DeviceName = device.DeviceName
			gpu.SubVendorName = device.SubVendorName
			gpu.SubDeviceName = device.SubDeviceName

			gpus = append(gpus, gpu)
		}
	}
	return gpus, nil
}
