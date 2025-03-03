package selector

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/types"
	"github.com/canonical/ml-snap-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

// TopStack finds the "best" stack to use. The current best choice is the smallest one to download to improve user experience.
func TopStack(scoredStacks []types.StackResult) (*types.StackResult, error) {
	if len(scoredStacks) == 0 {
		return nil, errors.New("no stacks found")
	}

	// Sort by size (low to high)
	sort.Slice(scoredStacks, func(i, j int) bool {
		return scoredStacks[i].Size < scoredStacks[j].Size
	})
	// return smallest
	return &scoredStacks[0], nil
}

func LoadStacksFromDir(stacksDir string) ([]types.Stack, error) {
	var stacks []types.Stack

	// Sanitise stack dir path
	if !strings.HasSuffix(stacksDir, "/") {
		stacksDir += "/"
	}

	// Iterate stacks
	files, err := os.ReadDir(stacksDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", stacksDir, err)
	}

	for _, file := range files {
		// Stacks dir should contain a dir per stack
		if !file.IsDir() {
			continue
		}

		data, err := os.ReadFile(stacksDir + file.Name() + "/stack.yaml")
		if err != nil {
			return nil, fmt.Errorf("%s: %s", stacksDir+file.Name(), err)
		}

		var currentStack types.Stack
		err = yaml.Unmarshal(data, &currentStack)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", stacksDir, err)
		}

		stacks = append(stacks, currentStack)
	}
	return stacks, nil
}

func FilterStacks(hardwareInfo types.HwInfo, stacks []types.Stack) ([]types.StackResult, error) {
	var scoredStacks []types.StackResult

	for _, currentStack := range stacks {
		compatible, notes, err := checkStack(hardwareInfo, currentStack)
		if err != nil {
			return nil, err
		}

		scoredStack := types.StackResult{
			Name:       currentStack.Name,
			Compatible: compatible,
			Notes:      notes,
		}

		stackSize, err := utils.StringToBytes(currentStack.DiskSpace)
		if err != nil {
			return nil, err
		}
		scoredStack.Size = stackSize

		scoredStacks = append(scoredStacks, scoredStack)
	}

	return scoredStacks, nil
}

func checkStack(hardwareInfo types.HwInfo, stack types.Stack) (bool, []string, error) {
	var notes []string
	var compatible = true

	// Enough memory
	requiredMemory, err := utils.StringToBytes(stack.Memory)
	if err != nil {
		return false, nil, err
	}

	if hardwareInfo.Memory == nil {
		notes = append(notes, "memory size not provided by hardware info")
		compatible = false
	} else {
		// Checking combination of ram and swap
		if hardwareInfo.Memory.TotalRam+hardwareInfo.Memory.TotalSwap < requiredMemory {
			notes = append(notes, "not enough memory")
			compatible = false
		}
	}

	// Enough disk space
	requiredDisk, err := utils.StringToBytes(stack.DiskSpace)
	if err != nil {
		return false, nil, err
	}
	if _, ok := hardwareInfo.Disk["/var/lib/snapd/snaps"]; !ok {
		notes = append(notes, "disk space not provided by hardware info")
		compatible = false
	} else {
		if hardwareInfo.Disk["/var/lib/snapd/snaps"].Avail < requiredDisk {
			notes = append(notes, "not enough free disk space")
			compatible = false
		}
	}

	// Devices
	// all
	allOfDevicesFound := 0
	for _, requiredDevice := range stack.Devices.All {
		switch requiredDevice.Type {
		case "cpu":
			if hardwareInfo.Cpus == nil {
				notes = append(notes, "cpu device is required but none found")
				compatible = false
			}
			result, reasons, err := checkCpus(requiredDevice, hardwareInfo.Cpus)
			if err != nil {
				return false, nil, err
			}
			if !result {
				notes = append(notes, reasons...)
				compatible = false
			}

		case "gpu":
			if len(hardwareInfo.Gpus) == 0 {
				notes = append(notes, "gpu device is required but none found")
				compatible = false
			}
			result, reasons, err := checkGpus(hardwareInfo.Gpus, requiredDevice)
			if err != nil {
				return false, nil, err
			}
			if !result {
				notes = append(notes, reasons...)
				compatible = false
			}
		}
		allOfDevicesFound++
	}

	if len(stack.Devices.All) > 0 && allOfDevicesFound != len(stack.Devices.All) {
		notes = append(notes, "all: could not find a required device")
		compatible = false
	}

	// any
	anyOfDevicesFound := 0
	for _, device := range stack.Devices.Any {
		switch device.Type {
		case "cpu":
			if hardwareInfo.Cpus == nil {
				continue
			}
			result, reasons, err := checkCpus(device, hardwareInfo.Cpus)
			if err != nil {
				return false, nil, err
			}
			notes = append(notes, reasons...)
			if !result {
				compatible = false
			}

		case "gpu":
			if hardwareInfo.Gpus == nil {
				continuereason
			}
			result, reasons, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return false, nil, err
			}
			if !result {
				notes = append(notes, reasons...)
				compatible = false
			}
		}
		anyOfDevicesFound++
	}

	// If any-of devices are defined, we need to find at least one
	if len(stack.Devices.Any) > 0 && anyOfDevicesFound == 0 {
		notes = append(notes, "any: could not find a required device")
		compatible = false
	}

	return compatible, notes, nil
}
