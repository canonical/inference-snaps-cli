package pci

import (
	"encoding/json"
	"fmt"
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

func TestParseLsCpu(t *testing.T) {
	devices, err := getDirectories("../../../test_data/devices")
	if err != nil {
		t.Fatal(err)
	}

	for _, device := range devices {
		lsPciFile := "../../../test_data/devices/" + device + "/lspci.txt"
		t.Run(device, func(t *testing.T) {
			_, err := os.Stat(lsPciFile)
			if err != nil {
				if os.IsNotExist(err) {
					// Device does not have lspci test data, skipping
					return
				} else {
					t.Fatal(err)
				}
			}

			lsPci, err := os.ReadFile(lsPciFile)
			if err != nil {
				t.Fatal(err)
			}

			pciDevices, err := ParseLsPci(string(lsPci), true)
			if err != nil {
				t.Fatal(err)
			}

			_, err = json.MarshalIndent(pciDevices, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
