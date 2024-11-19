package lspci

import (
	"encoding/json"
	"os"
	"testing"
)

var testFiles = []string{
	"test_data/xps13.txt",
}

func TestParseLsCpu(t *testing.T) {
	for _, lsPciFile := range testFiles {
		t.Run(lsPciFile, func(t *testing.T) {
			lsCpu, err := os.ReadFile(lsPciFile)
			if err != nil {
				t.Fatalf(err.Error())
			}

			cpuInfo, err := parseLsPci(lsCpu)
			if err != nil {
				t.Fatalf(err.Error())
			}

			jsonData, err := json.MarshalIndent(cpuInfo, "", "  ")
			if err != nil {
				t.Fatalf(err.Error())
			}

			t.Log(string(jsonData))
		})
	}
}
