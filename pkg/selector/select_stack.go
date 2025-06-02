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

func TopStack(scoredStacks []types.ScoredStack) (*types.ScoredStack, error) {
	var compatibleStacks []types.ScoredStack

	for _, stack := range scoredStacks {
		if stack.Score > 0 && stack.Grade == "stable" {
			compatibleStacks = append(compatibleStacks, stack)
		}
	}

	if len(compatibleStacks) == 0 {
		return nil, errors.New("no compatible stacks found")
	}

	// Sort by score (high to low) and return highest match
	sort.Slice(compatibleStacks, func(i, j int) bool {
		return compatibleStacks[i].Score > compatibleStacks[j].Score
	})

	// Top stack is highest score
	return &compatibleStacks[0], nil
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

func ScoreStacks(hardwareInfo types.HwInfo, stacks []types.Stack) ([]types.ScoredStack, error) {
	var scoredStacks []types.ScoredStack

	for _, currentStack := range stacks {
		score, err := checkStack(hardwareInfo, currentStack)

		scoredStack := types.ScoredStack{
			Stack:      currentStack,
			Score:      score,
			Compatible: true,
		}

		if score == 0 {
			scoredStack.Compatible = false
		}

		if err != nil {
			scoredStack.Notes = append(scoredStack.Notes, err.Error())
		}

		scoredStacks = append(scoredStacks, scoredStack)
	}

	return scoredStacks, nil
}

func checkStack(hardwareInfo types.HwInfo, stack types.Stack) (int, error) {
	stackScore := 0

	// Enough memory
	if stack.Memory != nil {
		requiredMemory, err := utils.StringToBytes(*stack.Memory)
		if err != nil {
			return 0, err
		}

		if hardwareInfo.Memory == nil {
			return 0, fmt.Errorf("no memory in hardware info")
		}

		// Checking combination of ram and swap
		if hardwareInfo.Memory.TotalRam+hardwareInfo.Memory.TotalSwap < requiredMemory {
			return 0, fmt.Errorf("not enough memory")
		}
		stackScore++
	}

	// Enough disk space
	if stack.DiskSpace != nil {
		requiredDisk, err := utils.StringToBytes(*stack.DiskSpace)
		if err != nil {
			return 0, err
		}
		if _, ok := hardwareInfo.Disk["/var/lib/snapd/snaps"]; !ok {
			return 0, fmt.Errorf("disk space not provided by hardware info")
		}
		if hardwareInfo.Disk["/var/lib/snapd/snaps"].Avail < requiredDisk {
			return 0, fmt.Errorf("not enough free disk space")
		}
		stackScore++
	}

	// Devices
	// all
	extraScore, err := checkDevicesAll(hardwareInfo, stack)
	if err != nil {
		return 0, fmt.Errorf("all: %v", err)
	}
	stackScore += extraScore

	// any
	extraScore, err = checkDevicesAny(hardwareInfo, stack)
	if err != nil {
		return 0, fmt.Errorf("any: %v", err)
	}
	stackScore += extraScore

	return stackScore, nil
}

func checkDevicesAll(hardwareInfo types.HwInfo, stack types.Stack) (int, error) {
	devicesFound := 0
	extraScore := 0

	for _, device := range stack.Devices.All {
		switch device.Type {
		case "cpu":
			if hardwareInfo.Cpus == nil {
				return 0, fmt.Errorf("cpu device is required but none found")
			}
			cpuScore, err := checkCpus(device, hardwareInfo.Cpus)
			if err != nil {
				return 0, fmt.Errorf("cpu: %v", err)
			}
			if cpuScore == 0 {
				return 0, fmt.Errorf("required cpu device not found")
			}
			extraScore += cpuScore
			devicesFound++

		case "gpu":
			if len(hardwareInfo.Gpus) == 0 {
				return 0, fmt.Errorf("gpu device is required but none found")
			}
			gpuScore, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return 0, fmt.Errorf("gpu: %v", err)
			}
			if gpuScore == 0 {
				return 0, fmt.Errorf("required gpu device not found")
			}
			extraScore += gpuScore
			devicesFound++

		default:
			// Devices without a type will land here - the type field should be empty
			if device.Type != "" {
				return 0, fmt.Errorf("stack device has an invalid type")
			}

			deviceScore, err := checkTypelessDevice(hardwareInfo, device)
			if err != nil {
				return 0, fmt.Errorf("error matching device: %v", err)
			}
			if deviceScore == 0 {
				return 0, fmt.Errorf("a required typeless device was not found")
			}
			extraScore += deviceScore
			devicesFound++
		}
	}

	if len(stack.Devices.All) > 0 && devicesFound != len(stack.Devices.All) {
		return 0, fmt.Errorf("could not find a required device")
	}

	return extraScore, nil
}

func checkDevicesAny(hardwareInfo types.HwInfo, stack types.Stack) (int, error) {
	devicesFound := 0
	extraScore := 0

	for _, device := range stack.Devices.Any {
		switch device.Type {
		case "cpu":
			if hardwareInfo.Cpus == nil {
				continue
			}
			cpuScore, err := checkCpus(device, hardwareInfo.Cpus)
			if err != nil {
				return 0, err
			}
			if cpuScore > 0 {
				devicesFound++
			}
			extraScore += cpuScore

		case "gpu":
			if hardwareInfo.Gpus == nil {
				continue
			}
			gpuScore, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return 0, err
			}
			if gpuScore > 0 {
				devicesFound++
			}
			extraScore += gpuScore

		default:
			// Devices without a type will land here - the type field should be empty
			if device.Type != "" {
				return 0, fmt.Errorf("stack device has an invalid type")
			}

			deviceScore, err := checkTypelessDevice(hardwareInfo, device)
			if err != nil {
				return 0, fmt.Errorf("error matching typeless device: %v", err)
			}
			if deviceScore > 0 {
				devicesFound++
			}
			extraScore += deviceScore
		}
	}

	// If any-of devices are defined, we need to find at least one
	if len(stack.Devices.Any) > 0 && devicesFound == 0 {
		return 0, fmt.Errorf("could not find a required device")
	}

	return extraScore, nil
}

func checkTypelessDevice(hardwareInfo types.HwInfo, device types.StackDevice) (int, error) {
	if device.Bus == nil {
		return 0, fmt.Errorf("stack devices without a type or bus are not supported")
	}
	bus := *device.Bus

	switch bus {
	case "pci":
		return checkPciPeripherals(hardwareInfo.PciPeripherals, device)
	case "usb":
		return 0, fmt.Errorf("usb devices not implemented")
	default:
		return 0, fmt.Errorf("unknown device bus")
	}
}
