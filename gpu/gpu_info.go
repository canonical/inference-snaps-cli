package gpu

import (
	"fmt"

	"github.com/canonical/hardware-info/lspci"
)

func Info() ([]Gpu, error) {
	var gpus []Gpu

	pciDevices, err := lspci.PciDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range pciDevices {
		// 00 01 - legacy VGA devices
		// 03 xx - display controllers
		if device.DeviceClass == 0x0001 || device.DeviceClass&0xFF00 == 0x0300 {
			var gpu Gpu
			gpu.VendorId = fmt.Sprintf("%04x", device.VendorId)
			gpu.DeviceId = fmt.Sprintf("%04x", device.DeviceId)
			if device.SubVendorId != nil {
				subVendorId := fmt.Sprintf("%04x", *device.SubVendorId)
				gpu.SubVendorId = &subVendorId
			}
			if device.SubDeviceId != nil {
				subDeviceId := fmt.Sprintf("%04x", *device.SubDeviceId)
				gpu.SubDeviceId = &subDeviceId
			}

			// TODO look up vram
			// TODO look up friendly names

			gpus = append(gpus, gpu)
		}
	}
	return gpus, nil
}
