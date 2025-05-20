package gpu

import (
	"os"
	"testing"
)

var clinfoFiles = []string{
	"../../../test_data/clinfo/ArcA580.json",
	"../../../test_data/clinfo/ArcB580.json",
	"../../../test_data/clinfo/NoDevices.json",
}

func TestParseClinfo(t *testing.T) {
	for _, clinfoFile := range clinfoFiles {
		t.Run(clinfoFile, func(t *testing.T) {
			clinfoJson, err := os.ReadFile(clinfoFile)
			if err != nil {
				t.Fatal(err)
			}
			clinfo, err := parseClinfoJson(clinfoJson)
			if err != nil {
				t.Fatal(err)
			}
			if len(clinfo.Devices) > 0 {
				if len(clinfo.Devices[0].Online) > 0 {
					t.Logf("Global memory: %d", clinfo.Devices[0].Online[0].ClDeviceGlobalMemSize)
				}
			}
		})
	}
}
