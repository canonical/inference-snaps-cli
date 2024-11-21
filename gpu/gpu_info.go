package gpu

import (
	"fmt"
	"github.com/canonical/hardware-info/lspci"
	"log"
)

func Info(friendlyNames bool) ([]Gpu, error) {
	pciDevices, err := lspci.PciDevices(friendlyNames)
	if err != nil {
		return nil, err
	}

	return pciGpus(pciDevices)
}

func pciGpus(pciDevices []lspci.PciDevice) ([]Gpu, error) {
	var gpus []Gpu

	for _, device := range pciDevices {
		// 00 01 - legacy VGA devices
		// 03 xx - display controllers
		if device.DeviceClass == 0x0001 || device.DeviceClass&0xFF00 == 0x0300 {
			var gpu Gpu
			gpu.VendorId = fmt.Sprintf("0x%04x", device.VendorId)
			gpu.DeviceId = fmt.Sprintf("0x%04x", device.DeviceId)
			if device.SubvendorId != nil {
				subVendorId := fmt.Sprintf("0x%04x", *device.SubvendorId)
				gpu.SubvendorId = &subVendorId
			}
			if device.SubdeviceId != nil {
				subDeviceId := fmt.Sprintf("0x%04x", *device.SubdeviceId)
				gpu.SubdeviceId = &subDeviceId
			}

			gpu.VendorName = device.VendorName
			gpu.DeviceName = device.DeviceName
			gpu.SubvendorName = device.SubvendorName
			gpu.SubdeviceName = device.SubdeviceName
			gpu.Properties = make(map[string]interface{})
			vendorSpecificInfo(&gpu, device)

			gpus = append(gpus, gpu)
		}
	}
	return gpus, nil
}

func vendorSpecificInfo(gpu *Gpu, pciDevice lspci.PciDevice) {

	switch pciDevice.VendorId {
	case 0x1002: // AMD
		vram, err := lookUpAmdVram(pciDevice)
		if err != nil {
			log.Printf("Error looking up AMD vRAM: %v", err)
		} else {
			gpu.Properties["vram"] = vram
		}

	case 0x10de: // NVIDIA
		vram, err := lookUpNvidiaVram(pciDevice)
		if err != nil {
			log.Printf("Error looking up NVIDIA vRAM: %v", err)
		} else {
			gpu.Properties["vram"] = vram
		}

		nvCompCap, err := computeCapability(pciDevice)
		if err != nil {
			log.Printf("Error looking up NVIDIA compute capability: %v", err)
		} else {
			gpu.Properties["compute_capability"] = nvCompCap
		}

	case 0x8086: // Intel
		log.Println("Not implemented")

	default:
		log.Println("Unknown GPU Vendor")
	}
}
