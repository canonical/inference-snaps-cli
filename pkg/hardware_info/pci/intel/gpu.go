package intel

import (
	"fmt"
	"os"
	"strconv"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func gpuProperties(pciDevice types.PciDevice) map[string]string {
	properties := make(map[string]string)

	vRamVal, err := vRam(pciDevice)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Intel: error looking up vRAM: %v", err)
	}
	if vRamVal != nil {
		properties["vram"] = strconv.FormatUint(*vRamVal, 10)
	}

	return properties
}

func vRam(device types.PciDevice) (*uint64, error) {
	return nil, fmt.Errorf("not implemented")
}
