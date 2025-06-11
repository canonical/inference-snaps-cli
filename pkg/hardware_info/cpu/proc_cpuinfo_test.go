package cpu

import (
	"os"
	"testing"
)

var procCpuInfoTestFiles = []string{
	"../../../test_data/cpuinfo/ampere-altra.txt",
	"../../../test_data/cpuinfo/ampere-one-m-banshee-12.txt",
	"../../../test_data/cpuinfo/ampere-one-siryn.txt",
	"../../../test_data/cpuinfo/ampere-one-x-banshee-8.txt",
	"../../../test_data/cpuinfo/raspberry-pi-5.txt",
}

func TestParseProcCpuInfo(t *testing.T) {

	for _, procCpuInfoFile := range procCpuInfoTestFiles {
		t.Run(procCpuInfoFile, func(t *testing.T) {
			procCpuInfoBytes, err := os.ReadFile(procCpuInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			parsed, err := parseProcCpuInfo(procCpuInfoBytes)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("parsed cpu info: %+v", parsed)

		})
	}
}
