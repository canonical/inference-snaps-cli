package gpu

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func intelVram(device pci.PciDevice) (*uint64, error) {
	/*
			For GPU vRAM information, I was able to see it with – clinfo. Grep for “Global memory size” and/or “Max memory allocation”
			After installing necessary drivers for GPU, NPU, you can use OV APIs to see available devices and their properties include VRAM

		clinfo --json
		CL_DEVICE_GLOBAL_MEM_SIZE
	*/
	command := exec.Command("clinfo", "--json")
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
			vram := clinfo.Devices[0].Online[0].ClDeviceGlobalMemSize
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
