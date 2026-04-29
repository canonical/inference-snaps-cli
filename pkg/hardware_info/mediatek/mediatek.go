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

func Info() ([]types.DetectedDevice, error) {
	compatibleValues, err := readCompatibleValues(deviceTreeRoot)
	if err != nil {
		return nil, fmt.Errorf("detecting mediatek apusys: %v", err)
	}

	devices := DetectFromCompatibles(compatibleValues)
	return devices, nil
}

func DetectFromCompatibles(compatibleValues []string) []types.DetectedDevice {
	platSet := make(map[string]struct{})

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
	}

	if len(platSet) == 0 {
		return nil
	}

	platformNames := make([]string, 0, len(platSet))
	for name := range platSet {
		platformNames = append(platformNames, name)
	}
	sort.Strings(platformNames)

	devices := make([]types.DetectedDevice, 0, len(platformNames))
	for _, name := range platformNames {
		devices = append(devices, types.DetectedDevice{
			Type: "npu",
			Bus:  "mdla",
			Metadata: &types.DeviceMetadata{
				VendorName:  "mediatek",
				ProductName: name,
			},
		})
	}

	return devices
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
