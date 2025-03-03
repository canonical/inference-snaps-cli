package selector

import (
	"fmt"
	"slices"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func checkCpus(stackDevice types.StackDevice, cpus []types.CpuInfo) (bool, []string, error) {
	var notes []string

	for _, cpu := range cpus {
		// First matching CPU returns true
		result, reason, err := cpuMatchesStack(cpu, stackDevice)
		if err != nil {
			return false, nil, err
		}
		notes = append(notes, reason)
		if result {
			return true, notes, nil
		}
	}
	// If we get here, we checked all the CPUs and none were matches
	return false, notes, nil
}

func cpuMatchesStack(cpu types.CpuInfo, stackDevice types.StackDevice) (bool, string, error) {
	// Vendor
	if stackDevice.VendorId != nil {
		if *stackDevice.VendorId != cpu.VendorId {
			return false, "cpu vendor mismatch", nil
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
			return false, fmt.Sprintf("cpu flag %s not found", flag), nil
		}
	}

	// Passed all filters, so it is a match
	return true, "cpu passed all checks", nil
}
