package intel

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/canonical/inference-snaps-cli/pkg/hardware_info/pci/utils"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/jpm-canonical/go-opencl/cl"
)

const clInfoTimeout = 10 * time.Second

func gpuProperties(pciDevice types.PciDevice) (map[string]string, error) {
	properties := make(map[string]string)

	vRamVal, err := VRamOpenCl(pciDevice)
	if err != nil {
		return nil, fmt.Errorf("error looking up vRAM: %v", err)
	}
	if vRamVal != nil {
		properties["vram"] = strconv.FormatUint(*vRamVal, 10)
	}

	return properties, nil
}

func vRam(device types.PciDevice) (*uint64, error) {
	/*
		For GPU vRAM information use clinfo. Grep for "Global memory size" and/or "Max memory allocation".
		After installing necessary drivers for GPU, NPU, you can also use OpenVino APIs to see available devices and their properties, including VRAM.
		`clinfo --json` reports a field `CL_DEVICE_GLOBAL_MEM_SIZE` which corresponds to the installed hardware's vRAM.
	*/
	ctx := context.Background()
	cmdContext, cancel := context.WithTimeout(ctx, clInfoTimeout)
	defer cancel()

	command := exec.CommandContext(cmdContext, "clinfo", "--json")

	// Set process group and kill the entire process tree on cancel
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	command.Cancel = func() error {
		return syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
	}

	data, err := command.Output()
	if err != nil {
		return nil, err
	}

	clinfo, err := parseClinfoJson(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse clinfo json: %w", err)
	}
	if len(clinfo.Devices) == 0 {
		return nil, fmt.Errorf("clinfo: no devices found")
	}
	if len(clinfo.Devices[0].Online) == 0 {
		return nil, fmt.Errorf("clinfo: no online devices found")
	}

	var vramValue *uint64 = nil
	// Search for the device with a matching PCI address
	for _, clInfoDevice := range clinfo.Devices[0].Online {
		if strings.Contains(clInfoDevice.ClDevicePciBusInfoKhr, device.Slot) {
			vram := clInfoDevice.ClDeviceGlobalMemSize
			vramValue = &vram
		}
	}
	return vramValue, nil
}

func parseClinfoJson(clinfoJson []byte) (types.Clinfo, error) {
	clinfo := types.Clinfo{}
	err := json.Unmarshal(clinfoJson, &clinfo)
	return clinfo, err
}

func VRamOpenCl(device types.PciDevice) (*uint64, error) {
	for _, platform := range cl.Platforms {
		for _, clDevice := range platform.Devices {
			// Print PCI bus info if cl_khr_pci_bus_info is available
			ext := clDevice.Property(cl.DEVICE_EXTENSIONS)

			extStr, ok := ext.(string)
			if !ok || strings.TrimSpace(extStr) == "" {
				return nil, fmt.Errorf("failed to look up device extensions")
			}

			if !strings.Contains(extStr, "cl_khr_pci_bus_info") {
				return nil, fmt.Errorf("device does not support pci bus info")
			}

			busInfo, err := clDevice.BusInfo()
			if err != nil {
				fmt.Println("error getting bus info:", err)
			}

			slotAddress, err := utils.ParsePCISlot(device.Slot)
			if err != nil {
				fmt.Println("error parsing pci slot address:", err)
			}

			if busInfo.Domain == slotAddress.Domain &&
				busInfo.Bus == slotAddress.Bus &&
				busInfo.Device == slotAddress.Device &&
				busInfo.Function == slotAddress.Function {
				globalMemory, err := clDevice.GlobalMemory()
				if err != nil {
					fmt.Println("error getting global memory:", err)
				}
				return &globalMemory, nil
			}
		}
	}
	return nil, fmt.Errorf("device not found in clinfo")
}
