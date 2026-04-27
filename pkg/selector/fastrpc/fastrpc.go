package fastrpc

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

const defaultNpuNodeGlob = "/dev/fastrpc-cdsp*"

func Match(manifestDevice engines.Device, hostDevices []types.DetectedDevice) (int, []string) {
	nodeGlob := defaultNpuNodeGlob
	if manifestDevice.NodeGlob != nil {
		nodeGlob = *manifestDevice.NodeGlob
	}

	for _, hostDevice := range hostDevices {
		if hostDevice.Bus != "fastrpc" {
			continue
		}
		if !deviceTypeMatch(manifestDevice.Type, hostDevice.Type) {
			continue
		}
		if hostDevice.Metadata == nil || hostDevice.Metadata.ProductName == "" {
			continue
		}

		matched, err := filepath.Match(nodeGlob, hostDevice.Metadata.ProductName)
		if err != nil {
			return 0, []string{fmt.Sprintf("invalid node-glob %q: %v", nodeGlob, err)}
		}
		if !matched {
			continue
		}

		score := weights.PciDevice + weights.PciDeviceType
		for _, connection := range manifestDevice.SnapConnections {
			if testing.Testing() {
				continue
			}
			connected, err := snapctl.IsConnected(connection).Run()
			if err != nil {
				return 0, []string{fmt.Sprintf("checking snap connection %q: %v", connection, err)}
			}
			if !connected {
				return 0, []string{fmt.Sprintf("%q is not connected", connection)}
			}
		}

		return score, nil
	}

	return 0, []string{fmt.Sprintf("device node matching %q not found", nodeGlob)}
}

func deviceTypeMatch(manifestType, hostType string) bool {
	if manifestType == "" {
		return true
	}

	manifest := strings.ToLower(strings.TrimSpace(manifestType))
	host := strings.ToLower(strings.TrimSpace(hostType))

	if manifest == host {
		return true
	}

	if manifest == "npu" && strings.HasPrefix(host, "npu") {
		return true
	}

	return false
}
