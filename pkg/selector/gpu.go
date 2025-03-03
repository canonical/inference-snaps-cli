package selector

import (
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/types"
	"github.com/canonical/ml-snap-utils/pkg/utils"
)

func checkGpus(gpus []types.Gpu, stackDevice types.StackDevice) (bool, []string, error) {
	var notes []string
	for _, systemGpu := range gpus {
		result, reason, err := gpuMatchesStack(systemGpu, stackDevice)
		if err != nil {
			return false, nil, err
		}
		notes = append(notes, reason)
		if result {
			// At the first matching GPU we stop and return
			return true, notes, nil
		}
	}
	// If we get here, we checked all the GPUs and none were matches
	return false, notes, nil
}

// gpuMatchesStack checks if the GPU matches what is required by the stack definition.
// This is done as a filter, based on the fields in the stack definition.
// If the GPU from the hardware info passes all these filters, the GPU is a match.
func gpuMatchesStack(gpu types.Gpu, stackDevice types.StackDevice) (bool, string, error) {

	// If the stack has a Vendor ID requirement, check if the GPU's vendor matches
	// Vendor IDs are hex number strings, so do a case-insensitive compare
	if stackDevice.VendorId != nil && !strings.EqualFold(*stackDevice.VendorId, gpu.VendorId) {
		return false, "gpu vendor mismatch", nil
	}

	// If stack has a vram requirement, check if GPU has enough
	if stackDevice.MinimumVram != nil {
		vramRequired, err := utils.StringToBytes(*stackDevice.MinimumVram)
		if err != nil {
			return false, "stack gpu vram invalid format", err
		}
		if gpu.VRam != nil {
			if *gpu.VRam < vramRequired {
				// Not enough vram
				return false, "gpu not enough vram", nil
			}
		} else {
			// Hardware Info does not list available vram
			return false, "gpu vram not reported", nil
		}
	}

	// TODO model id, compute capabilities

	// If we get here, all the filters have passed and the GPU is a match for this stack device
	return true, "gpu passed all checks", nil
}
