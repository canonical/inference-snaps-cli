package cpu

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGetLocalLscpu(t *testing.T) {
	output, err := GetLocalLscpu()
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println(string(output))

	//var lscpuJson LscpuContainer
	//_ = json.Unmarshal(output, &lscpuJson)
	//
	//jsonData, _ := json.MarshalIndent(lscpuJson, "", "  ")
	//fmt.Println(string(jsonData))
}

func TestParseLscpu(t *testing.T) {
	localLscpu, _ := GetLocalLscpu()
	cpuInfo := ParseLscpu(localLscpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLscpuR430(t *testing.T) {
	lscpu, _ := os.ReadFile("test_data/dell-r430-lscpu.json")
	cpuInfo := ParseLscpu(lscpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLscpuDl380p(t *testing.T) {
	lscpu, _ := os.ReadFile("test_data/hp-dl380p-gen8-lscpu.json")
	cpuInfo := ParseLscpu(lscpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLscpuRpi5(t *testing.T) {
	lscpu, _ := os.ReadFile("test_data/rpi5-lscpu.json")
	cpuInfo := ParseLscpu(lscpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLscpuMtG1200(t *testing.T) {
	lscpu, _ := os.ReadFile("test_data/mediatek-genio-1200-lscpu.json")
	cpuInfo := ParseLscpu(lscpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLscpuMtG350(t *testing.T) {
	lscpu, _ := os.ReadFile("test_data/mediatek-g350-lscpu.json")
	cpuInfo := ParseLscpu(lscpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}
