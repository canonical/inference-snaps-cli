package selector

import (
	"slices"

	types2 "github.com/canonical/ml-snap-utils/pkg/hardware_info/types"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func checkCpus(stackDevice types.StackDevice, cpu types2.CpuInfo) (int, error) {
	cpuScore := 0

	// Vendor
	if stackDevice.VendorId != nil {
		if *stackDevice.VendorId == cpu.Vendor {
			cpuScore += WeightCpuVendor // vendor matched
		} else {
			return 0, nil
		}
	}

	// TODO
	// architecture
	// cpu count

	for _, cpuModel := range cpu.Models {
		modelScore, err := checkCpuModel(cpuModel, stackDevice)
		if err != nil {
			return 0, err
		}
		if modelScore > 0 {
			// At the first matching CPU model we stop and return
			return cpuScore + modelScore, nil
		}
	}

	// If we get here, we checked all the CPU models and none were matches
	return 0, nil
}

// Apply the same "filter" logic as we have for the GPUs. See checkGpus() and checkGpu().
func checkCpuModel(cpuModel types2.Model, stackDevice types.StackDevice) (int, error) {
	// Each CPU that matches increases the score
	score := WeightCpu

	// Flags
	for _, flag := range stackDevice.Flags {
		if !slices.Contains(cpuModel.Flags, flag) {
			return 0, nil
		}
		score += WeightCpuFlag
	}

	// TODO
	// Family
	// CpuModel

	// If we get here, all the filters passed and the device is a match
	return score, nil
}
