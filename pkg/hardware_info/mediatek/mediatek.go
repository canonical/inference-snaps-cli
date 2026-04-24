package mediatek

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

var deviceTreeRoot = "/sys/firmware/devicetree/base"

var socToPlatform = map[string]string{
	"mt8195": "genio-1200",
	"mt8188": "genio-510-700",
	"mt8189": "genio-520-720",
}

func Info() ([]types.PlatformInfo, []types.DetectedDevice, error) {
	compatibleValues, err := readCompatibleValues(deviceTreeRoot)
	if err != nil {
		return nil, nil, fmt.Errorf("detecting mediatek apusys: %v", err)
	}

	platforms, devices := DetectFromCompatibles(compatibleValues)
	return platforms, devices, nil
}

func DetectFromCompatibles(compatibleValues []string) ([]types.PlatformInfo, []types.DetectedDevice) {
	platSet := make(map[string]struct{})
	matchSet := make(map[string]struct{})

	for _, value := range compatibleValues {
		soc, ok := extractAPUSoC(value)
		if !ok {
			continue
		}

		platformName, found := socToPlatform[soc]
		if !found {
			platformName = "mediatek-" + soc
		}

		platSet[platformName] = struct{}{}
		matchSet[value] = struct{}{}
	}

	if len(matchSet) == 0 {
		return nil, nil
	}

	platformNames := make([]string, 0, len(platSet))
	for name := range platSet {
		platformNames = append(platformNames, name)
	}
	sort.Strings(platformNames)

	platforms := make([]types.PlatformInfo, 0, len(platformNames))
	for _, name := range platformNames {
		platforms = append(platforms, types.PlatformInfo{
			Vendor: "mediatek",
			Name:   name,
		})
	}

	matches := make([]string, 0, len(matchSet))
	for value := range matchSet {
		matches = append(matches, value)
	}
	sort.Strings(matches)

	devices := []types.DetectedDevice{{
		Type:  "npu",
		Bus:   "mdla",
		Nodes: matches,
	}}

	return platforms, devices
}

func readCompatibleValues(root string) ([]string, error) {
	var values []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "compatible" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		entries := parseCompatibleData(data)
		values = append(values, entries...)
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return values, nil
}

func parseCompatibleData(data []byte) []string {
	parts := strings.Split(string(data), "\x00")
	values := make([]string, 0, len(parts))

	for _, part := range parts {
		v := strings.TrimSpace(part)
		if v == "" {
			continue
		}
		values = append(values, v)
	}

	return values
}

func extractAPUSoC(compatible string) (string, bool) {
	if !strings.HasPrefix(compatible, "mediatek,") || !strings.HasSuffix(compatible, "-apusys_rv") {
		return "", false
	}

	body := strings.TrimPrefix(compatible, "mediatek,")
	soc := strings.TrimSuffix(body, "-apusys_rv")
	if soc == "" || strings.Contains(soc, ",") {
		return "", false
	}

	return soc, true
}
