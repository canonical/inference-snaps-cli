package intel

import (
	"os"
	"strconv"

	"github.com/canonical/famous-models-cli/pkg/types"
	"github.com/canonical/go-snapctl"
)

func npuProperties(pciDevice types.PciDevice) (map[string]string, error) {
	properties := make(map[string]string)

	// We need to verify that the NPU driver is installed and available
	// To do this, we check if the correct Snap interface is connected
	// The method used is only possible from inside a snap
	if os.Getenv("SNAP") != "" {
		hasDriver, err := snapctl.IsConnected("npu-libs").Run()
		if err != nil {
			return nil, err
		}
		properties["has_driver"] = strconv.FormatBool(hasDriver)
	}

	return properties, nil
}
