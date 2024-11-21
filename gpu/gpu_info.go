package gpu

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

			vram, err := lookupVram(device)
			if err != nil {
				// If vram lookup fails, we only log it, as it is not fatal
				log.Println(err.Error())
			} else {
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

func lookupVram(device lspci.PciDevice) (uint64, error) {

	switch device.VendorId {
	case 0x1002: // AMD
		return lookupAmdVram(device)

	case 0x10de: // NVIDIA
		return lookupNvidiaVram(device)

	case 0x8086: // Intel
		return 0, errors.New("vram detection for Intel not implemented")
	}

	// TODO perhaps we can use eglinfo as a fallback
	return 0, errors.New(fmt.Sprintf("unable to detect vram for vendor %04x", device.VendorId))
}

func lookupAmdVram(device lspci.PciDevice) (uint64, error) {
	/*
		AMD vram is listed under /sys/bus/pci/devices/${pci_slot}/mem_info_vram_total

		ubuntu@u-HP-EliteBook-845-G8-Notebook-PC:~$ cat /sys/bus/pci/devices/0000\:04\:00.0/mem_info_
		mem_info_gtt_total       mem_info_vis_vram_total  mem_info_vram_used
		mem_info_gtt_used        mem_info_vis_vram_used   mem_info_vram_vendor
		mem_info_preempt_used    mem_info_vram_total

		ubuntu@u-HP-EliteBook-845-G8-Notebook-PC:~$ cat /sys/bus/pci/devices/0000\:04\:00.0/mem_info_vram_total
		536870912
	*/
	data, err := os.ReadFile("/sys/bus/pci/devices/" + device.Slot + "/mem_info_vram_total")
	if err != nil {
		return 0, err
	} else {
		dataStr := string(data)
		dataStr = strings.TrimSpace(dataStr) // value in file ends in \n
		return strconv.ParseUint(dataStr, 10, 64)
	}
}

func lookupNvidiaVram(device lspci.PciDevice) (uint64, error) {
	/*
		Nvidia: LANG=C nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits

		$ nvidia-smi --id=00000000:01:00.0 --query-gpu=memory.total --format=csv,noheader
		4096 MiB
		$ nvidia-smi --id=00000000:02:00.0 --query-gpu=memory.total --format=csv,noheader
		No devices were found
	*/
	command := exec.Command("nvidia-smi", "--id="+device.Slot, "--query-gpu=memory.total", "--format=csv,noheader")
	command.Env = os.Environ()
	command.Env = append(command.Env, "LANG=C")
	data, err := command.Output()
	if err != nil {
		return 0, err
	} else {
		dataStr := string(data)
		dataStr = strings.TrimSpace(dataStr) // value ends in \n
		valueStr, unit, hasUnit := strings.Cut(dataStr, " ")
		vramValue, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			return 0, err
		}

		if hasUnit {
			switch unit {
			case "KiB":
				vramValue = vramValue * 1024
			case "MiB":
				vramValue = vramValue * 1024 * 1024
			case "GiB":
				vramValue = vramValue * 1024 * 1024 * 1024
			}
		}

		return vramValue, nil
	}
}
