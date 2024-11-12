package cpu

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestGetHostLsCpu(t *testing.T) {
	hostLsCpu, err := GetHostLsCpu()
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println(string(hostLsCpu))
}

func TestParseHostLsCpu(t *testing.T) {
	hostLsCpu, err := GetHostLsCpu()
	if err != nil {
		t.Fatalf(err.Error())
	}

	cpuInfo, err := ParseLsCpu(hostLsCpu)
	if err != nil {
		t.Fatalf(err.Error())
	}

	jsonData, err := json.MarshalIndent(cpuInfo, "", "  ")
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(string(jsonData))
}

var testFiles = []string{
	"test_data/dell-r430-lscpu.json",
	"test_data/hp-dl380p-gen8-lscpu.json",
	"test_data/rpi5-lscpu.json",
	"test_data/mediatek-genio-1200-lscpu.json",
	"test_data/mediatek-g350-lscpu.json",
}

func TestParseLsCpu(t *testing.T) {
	for _, lsCpuFile := range testFiles {
		t.Run(lsCpuFile, func(t *testing.T) {
			lsCpu, err := os.ReadFile(lsCpuFile)
			if err != nil {
				t.Fatalf(err.Error())
			}

			cpuInfo, err := ParseLsCpu(lsCpu)
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
