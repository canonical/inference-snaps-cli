package selector

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/types"
	"github.com/canonical/ml-snap-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

func BestStack(compatibleStacks []types.StackResult) (*types.StackResult, error) {
	// Sort by score (high to low) and return best match
	sort.Slice(compatibleStacks, func(i, j int) bool {
		return compatibleStacks[i].Score > compatibleStacks[j].Score
	})

	// TODO find duplicate scores, use a different metric to choose one of them

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

func ScoreStacks(hardwareInfo types.HwInfo, stacks []types.Stack) ([]types.StackResult, error) {
	var scoredStacks []types.StackResult

	//tpuStacks := stacksForType(stacks, "tpu")
	gpuStacks := stacksForType(stacks, "gpu")
	cpuStacks := stacksForType(stacks, "cpu")

	// 1. TODO If hwinfo contains a TPU/NPU, check these stacks first

	// 2. Next check generic GPUs
	if len(hardwareInfo.Gpus) > 0 {
		log.Println("System has a GPU. Checking GPU stacks.")

		gpuStackMatches := 0
		for _, currentStack := range gpuStacks {
			score, err := checkStack(hardwareInfo, currentStack)
			if score > 0 {
				gpuStackMatches++
			}

			scoredStack := types.StackResult{
				Name:           currentStack.Name,
				Components:     currentStack.Components,
				Configurations: currentStack.Configurations,
				Score:          score,
			}

			if err != nil {
				scoredStack.Comment = err.Error()
			}

			scoredStacks = append(scoredStacks, scoredStack)
		}

		log.Printf("Found %d matching GPU stacks", gpuStackMatches)
		if gpuStackMatches > 0 {
			return scoredStacks, nil
		}
	}

	// 3. If no stacks have been found and returned, check CPU stacks
	log.Println("Checking CPU stacks.")

	cpuStackMatches := 0
	for _, currentStack := range cpuStacks {
		score, err := checkStack(hardwareInfo, currentStack)
		if score > 0 {
			cpuStackMatches++
		}

		foundStack := types.StackResult{
			Name:           currentStack.Name,
			Components:     currentStack.Components,
			Configurations: currentStack.Configurations,
			Score:          score,
		}

		if err != nil {
			foundStack.Comment = err.Error()
		}

		scoredStacks = append(scoredStacks, foundStack)
	}
	log.Printf("Found %d matching CPU stacks", cpuStackMatches)

	if len(scoredStacks) > 0 {
		return scoredStacks, nil
	}

	// If none found, return err
	return nil, fmt.Errorf("no stack found matching this hardware")
}

func stacksForType(allStacks []types.Stack, deviceType string) []types.Stack {
	var stacks []types.Stack
iterateStacks:
	for _, stack := range allStacks {
		for _, device := range stack.Devices.All {
			if device.Type == deviceType {
				stacks = append(stacks, stack)
				continue iterateStacks
			}
		}
		for _, device := range stack.Devices.Any {
			if device.Type == deviceType {
				stacks = append(stacks, stack)
				continue iterateStacks
			}
		}
	}
	return stacks
}

func checkStack(hardwareInfo types.HwInfo, stack types.Stack) (float64, error) {
	stackScore := 0.0

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
		if hardwareInfo.Memory.RamTotal+hardwareInfo.Memory.SwapTotal < requiredMemory {
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
	allOfDevicesFound := 0
	for _, device := range stack.Devices.All {
		switch device.Type {
		case "cpu":
			if hardwareInfo.Cpu == nil {
				return 0, fmt.Errorf("cpu device is required but none found")
			}
			cpuScore, err := checkCpus(device, *hardwareInfo.Cpu)
			if err != nil {
				return 0, err
			}
			if cpuScore == 0 {
				return 0, fmt.Errorf("required cpu device not found")
			}
			stackScore += cpuScore
			allOfDevicesFound++

		case "gpu":
			if len(hardwareInfo.Gpus) == 0 {
				return 0, fmt.Errorf("gpu device is required but none found")
			}
			gpuScore, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return 0, err
			}
			if gpuScore == 0 {
				return 0, fmt.Errorf("required gpu device not found")
			}
			stackScore += gpuScore
			allOfDevicesFound++
		}
	}

	if len(stack.Devices.All) > 0 && allOfDevicesFound != len(stack.Devices.All) {
		return 0, fmt.Errorf("all: could not find a required device")
	}

	// any
	anyOfDevicesFound := 0
	for _, device := range stack.Devices.Any {
		switch device.Type {
		case "cpu":
			if hardwareInfo.Cpu == nil {
				continue
			}
			cpuScore, err := checkCpus(device, *hardwareInfo.Cpu)
			if err != nil {
				return 0, err
			}
			if cpuScore > 0 {
				anyOfDevicesFound++
			}
			stackScore += cpuScore

		case "gpu":
			if hardwareInfo.Gpus == nil {
				continue
			}
			gpuScore, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return 0, err
			}
			if gpuScore > 0 {
				anyOfDevicesFound++
			}
			stackScore += gpuScore
		}
	}

	// If any-of devices are defined, we need to find at least one
	if len(stack.Devices.Any) > 0 && anyOfDevicesFound == 0 {
		return 0, fmt.Errorf("any: could not find a required device")
	}

	return stackScore, nil
}
