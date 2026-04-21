package qualcomm

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

var devRoot = "/dev"

func Info() ([]types.PlatformInfo, []types.DetectedDevice, error) {
	fastrpcNodes, err := globSorted(filepath.Join(devRoot, "fastrpc-*"))
	if err != nil {
		return nil, nil, fmt.Errorf("detecting fastrpc nodes: %v", err)
	}

	cdspNodes, err := globSorted(filepath.Join(devRoot, "fastrpc-cdsp*"))
	if err != nil {
		return nil, nil, fmt.Errorf("detecting cdsp fastrpc nodes: %v", err)
	}

	var platforms []types.PlatformInfo
	if len(fastrpcNodes) > 0 {
		platforms = append(platforms, types.PlatformInfo{
			Vendor: "qualcomm",
			Name:   "dragonwing",
		})
	}

	var devices []types.DetectedDevice
	if len(cdspNodes) > 0 {
		devices = append(devices, types.DetectedDevice{
			Type:  "npu",
			Bus:   "fastrpc",
			Nodes: cdspNodes,
		})
	}

	return platforms, devices, nil
}

func DetectFromNodes(nodes []string) ([]types.PlatformInfo, []types.DetectedDevice) {
	var fastrpcNodes []string
	var cdspNodes []string

	for _, node := range nodes {
		node = strings.TrimSpace(node)
		if node == "" {
			continue
		}

		base := filepath.Base(node)
		if strings.HasPrefix(base, "fastrpc-") {
			fastrpcNodes = append(fastrpcNodes, node)
		}
		if strings.HasPrefix(base, "fastrpc-cdsp") {
			cdspNodes = append(cdspNodes, node)
		}
	}

	sort.Strings(fastrpcNodes)
	sort.Strings(cdspNodes)

	var platforms []types.PlatformInfo
	if len(fastrpcNodes) > 0 {
		platforms = append(platforms, types.PlatformInfo{
			Vendor: "qualcomm",
			Name:   "dragonwing",
		})
	}

	var devices []types.DetectedDevice
	if len(cdspNodes) > 0 {
		devices = append(devices, types.DetectedDevice{
			Type:  "npu",
			Bus:   "fastrpc",
			Nodes: cdspNodes,
		})
	}

	return platforms, devices
}

func globSorted(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	return matches, nil
}
