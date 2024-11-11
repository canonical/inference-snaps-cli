package cpu

import "testing"

func TestGetInfo(t *testing.T) {
	GetProcCpuInfo()
}

func TestGetInfoDellR430(t *testing.T) {
	procCpuInfoPath = "test_data/dell-r430-cpuinfo"
	GetProcCpuInfo()
}

func TestGetInfoHpDl380p(t *testing.T) {
	procCpuInfoPath = "test_data/hp-dl380p-gen8-cpuinfo"
	GetProcCpuInfo()
}

func TestGetInfoXps13(t *testing.T) {
	procCpuInfoPath = "test_data/xps13-gen10-cpuinfo"
	GetProcCpuInfo()
}

func TestGetInfoG350(t *testing.T) {
	procCpuInfoPath = "test_data/mediatek-g350-cpuinfo"
	GetProcCpuInfo()
}

func TestGetInfoG1200(t *testing.T) {
	procCpuInfoPath = "test_data/mediatek-genio-1200-cpuinfo"
	GetProcCpuInfo()
}
