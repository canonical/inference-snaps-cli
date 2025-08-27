package types

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func getDirectories(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var directories []string
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}
	return directories, nil
}

func TestParseHwInfo(t *testing.T) {
	devices, err := getDirectories("../../test_data/devices")
	if err != nil {
		t.Fatal(err)
	}

	for _, device := range devices {
		hwInfoFile := "../../test_data/devices/" + device + "/hardware-info.json"
		t.Run(device, func(t *testing.T) {
			_, err := os.Stat(hwInfoFile)
			if err != nil {
				if os.IsNotExist(err) {
					// Device does not have hardware-info test data, skipping
					return
				} else {
					t.Fatal(err)
				}
			}

			file, err := os.Open(hwInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			data, err := io.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			var hardwareInfo HwInfo
			err = json.Unmarshal(data, &hardwareInfo)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
