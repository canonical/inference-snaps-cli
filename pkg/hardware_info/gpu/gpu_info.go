package gpu

import (
	"errors"
	"fmt"
	"log"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
	types2 "github.com/canonical/ml-snap-utils/pkg/types"
)

func Info(friendlyNames bool) ([]types2.Gpu, error) {
	pciDevices, err := pci.PciDevices(friendlyNames)
	if err != nil {
		return nil, err
	}

	return pciGpus(pciDevices)
}

func pciGpus(pciDevices []types2.Device) ([]types2.Gpu, error) {
	var gpus []types2.Gpu

	for _, device := range pciDevices {
		// 00 01 - legacy VGA devices
		// 03 xx - display controllers
		if device.DeviceClass == 0x0001 || device.DeviceClass&0xFF00 == 0x0300 {
			var gpu types2.Gpu
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

			vram, err := getVRam(device)
			if err != nil {
				log.Printf("Error getting VRAM info for GPU: %s", err)
			}
			gpu.VRam = vram

			computeCapability, err := getComputeCapability(device)
			if err != nil {
				log.Printf("Error getting compute capability for GPU: %s", err)
			}
			gpu.ComputeCapability = computeCapability

			gpus = append(gpus, gpu)
		}
	}
	return gpus, nil
}

func getVRam(pciDevice types2.Device) (*uint64, error) {
	switch pciDevice.VendorId {
	case 0x1002: // AMD
		return amdVram(pciDevice)
	case 0x10de: // NVIDIA
		return nvidiaVram(pciDevice)
	case 0x8086: // Intel
		return nil, errors.New("vram lookup for Intel GPU not implemented")
	default:
		return nil, errors.New("vram lookup for unknown GPU not implemented")
	}
}

func getComputeCapability(pciDevice types2.Device) (*string, error) {
	switch pciDevice.VendorId {
	case 0x1002: // AMD
		return nil, errors.New("compute capability lookup for AMD GPU not implemented")
	case 0x10de: // NVIDIA
		return nvidiaComputeCapability(pciDevice)
	case 0x8086: // Intel
		return nil, errors.New("compute capability lookup for Intel GPU not implemented")
	default:
		return nil, errors.New("compute capability lookup for unknown GPU not implemented")
	}
}
