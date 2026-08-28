package fastrpc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/lscompute/pkg/machine"
	lsfastrpc "github.com/canonical/lscompute/pkg/machine/device/fastrpc"
)

func Match(manifestDevice engines.Device, machineInfo *machine.MachineInfo) (int, []string) {
	if machineInfo == nil {
		return 0, []string{"no machine info provided"}
	}

	for _, device := range machineInfo.Devices {
		fastRPCDevice, ok := device.(lsfastrpc.Device)
		if !ok || fastRPCDevice.Bus != lsfastrpc.BusName {
			continue
		}
		if manifestDevice.Type != "" && manifestDevice.Type != "npu" {
			continue
		}
		// A domain-less manifest retains the original presence-only behavior.
		// Backends such as Hexagon HTP should require their domain explicitly.
		if manifestDevice.Domain != nil && !strings.EqualFold(strings.TrimSpace(*manifestDevice.Domain), string(fastRPCDevice.Domain)) {
			continue
		}

		for _, connection := range manifestDevice.SnapConnections {
			connected, err := checkSnapConnection(connection)
			if err != nil {
				return 0, []string{fmt.Sprintf("checking snap connection %q: %v", connection, err)}
			}
			if !connected {
				return 0, []string{fmt.Sprintf("%q is not connected", connection)}
			}
		}

		return weights.PciDevice + weights.PciDeviceType, nil
	}

	if manifestDevice.Domain != nil {
		return 0, []string{fmt.Sprintf("no fastrpc devices in %q domain on host system", *manifestDevice.Domain)}
	}
	return 0, []string{"no fastrpc devices on host system"}
}

func checkSnapConnection(connection string) (bool, error) {
	if testing.Testing() {
		// Tests do not necessarily run inside a snap.
		return true, nil
	}
	return snap.IsConnected(connection)
}
