package gpu

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

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

			vram, err := lookuUpVram(device)
			if err == nil {
				gpu.VRam = &vram
			}

			gpu.VendorName = device.VendorName
			gpu.DeviceName = device.DeviceName
			gpu.SubVendorName = device.SubVendorName
			gpu.SubDeviceName = device.SubDeviceName

			gpus = append(gpus, gpu)
		}
	}
	return gpus, nil
}

func lookuUpVram(device lspci.PciDevice) (uint64, error) {
	// AMD vram is listed under /sys/bus/pci/devices/${pci_slot}/mem_info_vram_total
	if device.VendorId == 0x1002 {
		/*
			ubuntu@u-HP-EliteBook-845-G8-Notebook-PC:~$ cat /sys/bus/pci/devices/0000\:04\:00.0/mem_info_
			mem_info_gtt_total       mem_info_vis_vram_total  mem_info_vram_used
			mem_info_gtt_used        mem_info_vis_vram_used   mem_info_vram_vendor
			mem_info_preempt_used    mem_info_vram_total

			ubuntu@u-HP-EliteBook-845-G8-Notebook-PC:~$ cat /sys/bus/pci/devices/0000\:04\:00.0/mem_info_vram_total
			536870912
		*/
		data, err := os.ReadFile("/sys/bus/pci/devices/" + device.Slot + "/mem_info_vram_total")
		if err == nil {
			size, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return size, nil
			}
		} else {
			log.Println("Failed to look up AMD VRAM")
		}
	}

	// Nvidia: LANG=C nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits
	if device.VendorId == 0x10de {
		out, err := exec.Command("LANG=C", "nvidia-smi", "--query-gpu=memory.total", "--format=csv,noheader,nounits").Output()
		if err == nil {
			size, err := strconv.ParseUint(string(out), 10, 64)
			if err != nil {
				return size, nil
			}
		} else {
			log.Println("Failed to look up NVIDIA VRAM")
		}
	}

	return 0, errors.New("unable to detect vram")
}
