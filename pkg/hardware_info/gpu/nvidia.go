package gpu

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
	"github.com/canonical/ml-snap-utils/pkg/utils"
)

func lookUpNvidiaVram(device pci.Device) (uint64, error) {
	/*
		$ LANG=C nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits

		$ nvidia-smi --id=00000000:01:00.0 --query-gpu=memory.total --format=csv,noheader
		4096 MiB
		$ nvidia-smi --id=00000000:02:00.0 --query-gpu=memory.total --format=csv,noheader
		No devices were found
	*/

	nvidiaSmi, err := findNvidiaSmi()
	if err != nil {
		return 0, err
	}
	command := exec.Command(nvidiaSmi, "--id="+device.Slot, "--query-gpu=memory.total", "--format=csv,noheader")
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

func computeCapability(device pci.Device) (string, error) {
	// $ LANG=C nvidia-smi --id=00000000:01:00.0 --query-gpu=compute_cap --format=csv,noheader

	nvidiaSmi, err := findNvidiaSmi()
	if err != nil {
		return "", err
	}
	command := exec.Command(nvidiaSmi, "--id="+device.Slot, "--query-gpu=compute_cap", "--format=csv,noheader")
	command.Env = os.Environ()
	command.Env = append(command.Env, "LANG=C")
	data, err := command.Output()
	if err != nil {
		return "", err
	}

	ccValue := strings.TrimSpace(string(data))
	return ccValue, nil
}

func findNvidiaSmi() (string, error) {
	nvidiaSmiFallback := "nvidia-smi" // Fall back to find in PATH
	hostFs := "/var/lib/snapd/hostfs"
	nvidiaSmiHostFs := hostFs + "/usr/bin/nvidia-smi"
	nvidiaSmiTmp := "/tmp/nvidia-smi"

	// If the hostfs path exists, we are inside a snap
	infoHostFs, err := os.Stat(hostFs)
	// No error means the path exists
	if err == nil {
		// If it is not readable, the system-backup interface is missing
		if infoHostFs.Mode().Perm()&0444 == 0444 {
			log.Println("Can't access hostfs. Try connecting the system-backup interface.")
			return nvidiaSmiFallback, nil
		}

		/*
			If we are running inside a snap, and the snap has access to the host file system via the system-backup interface,
			nvidia-smi from the host can be accessed. This path is read-only without execution, so we need to copy it to
			another location first and then fix permissions.
		*/
		_, err = os.Stat(nvidiaSmiHostFs)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Not inside a snap, or nvidia-smi is not on the host. Fall through to fallback.
				return nvidiaSmiFallback, nil
			} else {
				// A different error occurred
				return nvidiaSmiFallback, err
			}
		} else {
			err = utils.CopyFile(nvidiaSmiHostFs, nvidiaSmiTmp)
			if err != nil {
				return nvidiaSmiFallback, err
			}
			err = os.Chmod(nvidiaSmiTmp, 0755)
			if err != nil {
				return nvidiaSmiFallback, err
			}
			return nvidiaSmiTmp, nil
		}
	}

	// Not inside a snap, always use the fallback
	return nvidiaSmiFallback, nil
}
