package renesas

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

var devRoot = "/dev"

func Info() ([]types.DetectedDevice, error) {
	nodes, err := globSorted(filepath.Join(devRoot, "drpai*"))
	if err != nil {
		return nil, fmt.Errorf("detecting drpai nodes: %v", err)
	}

	return DetectFromNodes(nodes), nil
}

func DetectFromNodes(nodes []string) []types.DetectedDevice {
	var devices []types.DetectedDevice

	for _, node := range nodes {
		node = strings.TrimSpace(node)
		if node == "" {
			continue
		}

		base := filepath.Base(node)
		if !strings.HasPrefix(base, "drpai") {
			continue
		}

		devices = append(devices, types.DetectedDevice{
			Type: "npu",
			Bus:  "drpai",
			Metadata: &types.DeviceMetadata{
				VendorName:  "renesas",
				ProductName: node,
			},
		})
	}

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Metadata.ProductName < devices[j].Metadata.ProductName
	})

	return devices
}

func globSorted(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	sort.Strings(matches)
	return matches, nil
}
