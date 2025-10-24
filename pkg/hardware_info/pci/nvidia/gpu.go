package nvidia

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

const nvidiaSmiTimeout = 60 * time.Second

func gpuProperties(pciDevice types.PciDevice) (map[string]string, error) {
	properties := make(map[string]string)

	vRamVal, err := vRam(pciDevice)
	if err != nil {
		return nil, fmt.Errorf("error looking up vRAM: %v", err)
	}
	if vRamVal != nil {
		properties["vram"] = strconv.FormatUint(*vRamVal, 10)
	}

	ccVal, err := computeCapability(pciDevice)
	if err != nil {
		return nil, fmt.Errorf("error looking up compute capability: %v", err)
	}
	if ccVal != nil {
		properties["compute-capability"] = *ccVal
	}

	return properties, nil
}

func vRam(device types.PciDevice) (*uint64, error) {
	/*
		Nvidia: LANG=C nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits

		$ nvidia-smi --id=00000000:01:00.0 --query-gpu=memory.total --format=csv,noheader
		4096 MiB
		$ nvidia-smi --id=00000000:02:00.0 --query-gpu=memory.total --format=csv,noheader
		No devices were found
	*/
	output, err := nvidiaSmi("--id="+device.Slot, "--query-gpu=memory.total", "--format=csv,noheader")
	if err != nil {
		return nil, fmt.Errorf("error executing nvidia-smi: %s", err)
	}

	valueStr, unit, hasUnit := strings.Cut(*output, " ")
	vramValue, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return nil, err
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

	return &vramValue, nil

}

func computeCapability(device types.PciDevice) (*string, error) {
	// nvidia-smi --query-gpu=compute_cap --format=csv,noheader
	output, err := nvidiaSmi("--id="+device.Slot, "--query-gpu=compute_cap", "--format=csv,noheader")
	if err != nil {
		return nil, fmt.Errorf("error executing nvidia-smi: %s", err)
	}

	return output, nil
}

func nvidiaSmi(args ...string) (*string, error) {
	cmd := exec.Command("nvidia-smi", args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=C")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	type cmdOutputStruct struct {
		output []byte
		err    error
	}
	cmdOutputChannel := make(chan cmdOutputStruct, 1)
	go func() {
		output, err := cmd.CombinedOutput()
		cmdOutputChannel <- cmdOutputStruct{output, err}
	}()

	var data []byte
	select {
	case <-time.After(nvidiaSmiTimeout):
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return nil, fmt.Errorf("nvidia-smi timed out and killed")
	case cmdOutput, ok := <-cmdOutputChannel:
		if !ok {
			return nil, fmt.Errorf("channel closed unexpectedly")
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

	strOutput := string(bytes.TrimSpace(data))
	return &strOutput, nil
}
