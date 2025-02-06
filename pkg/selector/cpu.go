package selector

import (
	"slices"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func checkCpus(stackDevice types.StackDevice, cpus []types.CpuInfo) (int, error) {
	cpusScore := 0

iterateCpus:
	for _, cpu := range cpus {
		cpuScore := WeightCpu

		// Vendor
		if stackDevice.VendorId != nil {
			if *stackDevice.VendorId == cpu.VendorId {
				cpuScore += WeightCpuVendor // vendor matched
			} else {
				continue
			}
		}

		// TODO
		// architecture
		// cpu count
		// Family
		// CpuModel

		// Flags
		for _, flag := range stackDevice.Flags {
			if !slices.Contains(cpu.Flags, flag) {
				continue iterateCpus
			}
			cpuScore += WeightCpuFlag
		}
		cpusScore += cpuScore
	}

	// If we get here, we checked all the CPU models and none were matches
	return cpusScore, nil
}
