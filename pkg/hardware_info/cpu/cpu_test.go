package cpu

import (
	"os"
	"testing"
)

func TestEnrichArm(t *testing.T) {
	lsCpu, err := os.ReadFile("../../../test_data/lscpu/ampere-altra.json")
	if err != nil {
		t.Fatal(err)
	}

	cpus, err := parseLsCpu(lsCpu)
	if err != nil {
		t.Fatal(err)
	}

	procCpuInfoData, err := os.ReadFile("../../../test_data/cpuinfo/ampere-altra.txt")
	if err != nil {
		t.Fatal(err)
	}

	procCpuInfoParsed, err := parseProcCpuInfo(procCpuInfoData)
	if err != nil {
		t.Fatal(err)
	}

	cpus, err = enrichCpus(cpus, procCpuInfoParsed)
	if err != nil {
		t.Fatal(err)
	}

	if cpus[0].ModelId != 3340 {
		t.Fatal("cpu model enrichment failed")
	}

	t.Logf("%+v", cpus)
}
