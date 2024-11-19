package lspci

import (
	"encoding/json"
	"os"
	"testing"
)

var testFiles = []string{
	"test_data/dell-precision-3660-c29399.txt",
	"test_data/dell-vostro153535-c30942.txt",
	"test_data/farshid-amd.txt",
	"test_data/hp-elitebook845-g8-notebook-pc-c30368.txt",
	"test_data/katryn.txt",
	"test_data/lana.txt",
	"test_data/magda.txt",
	"test_data/rpi5.txt",
	"test_data/xps13.txt",
}

func TestParseLsCpu(t *testing.T) {
	for _, lsPciFile := range testFiles {
		t.Run(lsPciFile, func(t *testing.T) {
			lsPci, err := os.ReadFile(lsPciFile)
			if err != nil {
				t.Fatalf(err.Error())
			}

			pciDevices, err := ParseLsPci(lsPci, true)
			if err != nil {
				t.Fatalf(err.Error())
			}

			jsonData, err := json.MarshalIndent(pciDevices, "", "  ")
			if err != nil {
				t.Fatalf(err.Error())
			}

			t.Log(string(jsonData))
		})
	}
}
