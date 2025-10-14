package cpu

import (
	"log"
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/constants"
)

var procCpuInfoTestFiles = map[string]string{
	"../../../test_data/machines/ampere-one-m-banshee-12/cpuinfo.txt":           constants.Arm64,
	"../../../test_data/machines/ampere-one-siryn/cpuinfo.txt":                  constants.Arm64,
	"../../../test_data/machines/ampere-one-x-banshee-8/cpuinfo.txt":            constants.Arm64,
	"../../../test_data/machines/hp-proliant-rl300-gen11-altra/cpuinfo.txt":     constants.Arm64,
	"../../../test_data/machines/hp-proliant-rl300-gen11-altra-max/cpuinfo.txt": constants.Arm64,
	"../../../test_data/machines/i7-2600k+arc-a580/cpuinfo.txt":                 constants.Amd64,
	"../../../test_data/machines/i7-10510U/cpuinfo.txt":                         constants.Amd64,
	"../../../test_data/machines/mustang/cpuinfo.txt":                           constants.Amd64,
	//"../../../test_data/machines/orange-pi-rv2/cpuinfo.txt":                     constants.Riscv64,
	"../../../test_data/machines/raspberry-pi-5/cpuinfo.txt":         constants.Arm64,
	"../../../test_data/machines/raspberry-pi-5+hailo-8/cpuinfo.txt": constants.Arm64,
	"../../../test_data/machines/xps13-7390/cpuinfo.txt":             constants.Amd64,
	"../../../test_data/machines/xps13-9350/cpuinfo.txt":             constants.Amd64,
}

func TestParseProcCpuInfo(t *testing.T) {

	for procCpuInfoFile, arch := range procCpuInfoTestFiles {
		t.Run(procCpuInfoFile, func(t *testing.T) {
			procCpuInfoBytes, err := os.ReadFile(procCpuInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			parsed, err := parseProcCpuInfo(string(procCpuInfoBytes), arch)
			if err != nil {
				t.Fatal(err)
			}

			for _, cpuInfo := range parsed {
				log.Printf("%+v", cpuInfo)
			}

		})
	}
}

func TestParseProcCpuInfoAmd64(t *testing.T) {
	cpuInfoData, err := os.ReadFile("../../../test_data/machines/xps13-7390/cpuinfo.txt")
	if err != nil {
		t.Fatal(err)
	}

	cpuInfos, err := parseProcCpuInfoAmd64(string(cpuInfoData))
	if err != nil {
		t.Fatal(err)
	}

	for _, cpuInfo := range cpuInfos {
		log.Printf("%+v", cpuInfo)
	}
}

func TestParseProcCpuInfoArm64(t *testing.T) {
	cpuInfoData, err := os.ReadFile("../../../test_data/machines/raspberry-pi-5/cpuinfo.txt")
	if err != nil {
		t.Fatal(err)
	}

	cpuInfos, err := parseProcCpuInfoArm64(string(cpuInfoData))
	if err != nil {
		t.Fatal(err)
	}

	for _, cpuInfo := range cpuInfos {
		log.Printf("%+v", cpuInfo)
	}
}
