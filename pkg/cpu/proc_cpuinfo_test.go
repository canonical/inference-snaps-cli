package cpu

import "testing"

func TestGetInfo(t *testing.T) {
	GetInfo()
}

func TestGetInfoDellR430(t *testing.T) {
	procCpuInfoPath = "test_data/dell-r430-cpuinfo"
	GetInfo()
}

func TestGetInfoHpDl380p(t *testing.T) {
	procCpuInfoPath = "test_data/hp-dl380p-gen8-cpuinfo"
	GetInfo()
}

func TestGetInfoXps13(t *testing.T) {
	procCpuInfoPath = "test_data/xps13-gen10-cpuinfo"
	GetInfo()
}

func TestGetInfoG350(t *testing.T) {
	procCpuInfoPath = "test_data/mediatek-g350-cpuinfo"
	GetInfo()
}

func TestGetInfoG1200(t *testing.T) {
	procCpuInfoPath = "test_data/mediatek-genio-1200-cpuinfo"
	GetInfo()
}
