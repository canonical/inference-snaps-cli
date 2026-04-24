package qualcomm

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

var devRoot = "/dev"

func Info() ([]types.DetectedDevice, error) {
	nodes, err := globSorted(filepath.Join(devRoot, "fastrpc-*"))
	if err != nil {
		return nil, fmt.Errorf("detecting fastrpc nodes: %v", err)
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
		if !strings.HasPrefix(base, "fastrpc-") {
			continue
		}

		devices = append(devices, types.DetectedDevice{
			Type: detectNpuType(base),
			Bus:  "fastrpc",
			PlatformInfo: &types.PlatformInfo{
				Vendor: "qualcomm",
				Name:   node,
				SoC:    "dragonwing",
			},
		})
	}

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].PlatformInfo.Name < devices[j].PlatformInfo.Name
	})

	return devices
}

func detectNpuType(nodeBaseName string) string {
	const prefix = "fastrpc-"
	if strings.HasPrefix(nodeBaseName, prefix) {
		suffix := strings.TrimPrefix(nodeBaseName, prefix)
		return fmt.Sprintf("NPU - %s", suffix)
	}
	return "NPU"
}

func globSorted(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	sort.Strings(matches)
	return matches, nil
}
