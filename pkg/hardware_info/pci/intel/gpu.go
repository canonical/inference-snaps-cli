package intel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

const clInfoTimeout = 30 * time.Second

func gpuProperties(pciDevice types.PciDevice) (map[string]string, error) {
	properties := make(map[string]string)

	vRamVal, err := vRam(pciDevice)
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

		We add a timeout to prevent clinfo from blocking the program flow.
		See https://jarv.org/posts/command-with-timeout/ for a guide on properly timing out a command in Go.
	*/
	command := exec.Command("clinfo", "--json")
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	type cmdOutputStruct struct {
		output []byte
		err    error
	}
	cmdOutputChannel := make(chan cmdOutputStruct, 1)
	go func() {
		output, err := command.CombinedOutput()
		cmdOutputChannel <- cmdOutputStruct{output, err}
	}()

	var data []byte
	select {
	case <-time.After(clInfoTimeout):
		syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
		return nil, fmt.Errorf("clinfo timed out and killed")
	case cmdOutput, ok := <-cmdOutputChannel:
		if !ok {
			return nil, fmt.Errorf("command channel closed unexpectedly")
		}
		if cmdOutput.err != nil {
			if len(cmdOutput.output) == 0 {
				return nil, cmdOutput.err
			} else {
				return nil, fmt.Errorf("%s: %s", cmdOutput.err, bytes.TrimSpace(cmdOutput.output))
			}
		}
		data = cmdOutput.output
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
