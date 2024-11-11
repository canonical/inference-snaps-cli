package cpu

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGetLocalLsCpu(t *testing.T) {
	hostLsCpu, err := GetHostLsCpu()
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println(string(hostLsCpu))
}

func TestParseLsCpu(t *testing.T) {
	hostLsCpu, _ := GetHostLsCpu()
	cpuInfo := ParseLsCpu(hostLsCpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLsCpuR430(t *testing.T) {
	lsCpu, _ := os.ReadFile("test_data/dell-r430-lscpu.json")
	cpuInfo := ParseLsCpu(lsCpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLsCpuDl380p(t *testing.T) {
	lsCpu, _ := os.ReadFile("test_data/hp-dl380p-gen8-lscpu.json")
	cpuInfo := ParseLsCpu(lsCpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLsCpuRpi5(t *testing.T) {
	lsCpu, _ := os.ReadFile("test_data/rpi5-lscpu.json")
	cpuInfo := ParseLsCpu(lsCpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLsCpuMtG1200(t *testing.T) {
	lsCpu, _ := os.ReadFile("test_data/mediatek-genio-1200-lscpu.json")
	cpuInfo := ParseLsCpu(lsCpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLsCpuMtG350(t *testing.T) {
	lsCpu, _ := os.ReadFile("test_data/mediatek-g350-lscpu.json")
	cpuInfo := ParseLsCpu(lsCpu)

	jsonData, _ := json.MarshalIndent(cpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}
