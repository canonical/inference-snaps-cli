package cpu

import (
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/pkg/constants"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

/*
Match takes a Device with type CPU, and checks if it matches any of the CPU models reported for the system.
A score, a string slice with reasons and an error are returned. If there is a matching CPU on the system, the score will be positive and the error will be nil.
If no CPU is found, the score will be zero and there will be one or more reasons for the mismatch. In case of a runtime error, the error value will be non-nil.
*/
func Match(device engines.Device, cpus []types.CpuInfo) (maxCpuScore int, deviceIssues []string) {
	maxCpuScore = 0

	if cpus == nil {
		deviceIssues = append(deviceIssues, "no cpu found on host system")
	}

	for i, cpu := range cpus {
		cpuScore, cpuIssues := CheckCpu(device, cpu)

		if len(cpuIssues) > 0 {
			if len(cpus) > 1 {
				for _, issue := range cpuIssues {
					deviceIssues = append(deviceIssues, fmt.Sprintf("cpu %d: %v", i, issue))
				}
			} else {
				deviceIssues = append(deviceIssues, cpuIssues...)
			}
		} else {
			if cpuScore > maxCpuScore {
				maxCpuScore = cpuScore
			}
		}
	}

	return
}

func CheckCpu(device engines.Device, cpu types.CpuInfo) (cpuScore int, issues []string) {
	cpuScore += weights.CpuDevice

	// architecture
	if device.Architecture != nil {
		if *device.Architecture == cpu.Architecture {
			// architecture matches - no additional weight
		} else {
			issues = append(issues, fmt.Sprintf("incorrect architecture %s", cpu.Architecture))
		}
	}

	/*
		amd64
	*/
	if cpu.Architecture == constants.Amd64 {
		// amd64 manufacturer ID
		if device.ManufacturerId != nil {
			if *device.ManufacturerId == cpu.ManufacturerId {
				cpuScore += weights.CpuVendor
			} else {
				issues = append(issues, fmt.Sprintf("manufacturer mismatch: %s", cpu.ManufacturerId))
			}
		}

		// amd64 flags
		for _, flag := range device.Flags {
			if slices.Contains(cpu.Flags, flag) {
				cpuScore += weights.CpuFlag
			} else {
				issues = append(issues, fmt.Sprintf("flag not available: %s", flag))
			}
		}
	}

	/*
		arm64
	*/
	if cpu.Architecture == constants.Arm64 {
		// arm64 implementer ID
		if device.ImplementerId != nil {
			if *device.ImplementerId == cpu.ImplementerId {
				cpuScore += weights.CpuVendor
			} else {
				issues = append(issues, fmt.Sprintf("implementer id mismatch: %x", cpu.ImplementerId))
			}
		}

		// arm64 part number
		if device.PartNumber != nil {
			if *device.PartNumber == cpu.PartNumber {
				cpuScore += weights.CpuModel
			} else {
				issues = append(issues, fmt.Sprintf("part number mismatch: %x", cpu.PartNumber))
			}
		}

		// arm64 features
		for _, feature := range device.Features {
			if slices.Contains(cpu.Features, feature) {
				cpuScore += weights.CpuFlag
			} else {
				issues = append(issues, fmt.Sprintf("feature not found: %s", feature))
			}
		}
	}

	if len(issues) > 0 {
		cpuScore = 0
	}

	return
}
