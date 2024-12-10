package selector

import (
	"testing"

	"github.com/canonical/hardware-info/hardware_info/cpu"
	"github.com/canonical/hardware-info/types"
)

func TestCheckCpu(t *testing.T) {
	vendorId := "GenuineIntel"
	stackDevice := types.StackDevice{
		Type:     "cpu",
		Bus:      nil,
		VendorId: &vendorId,
	}

	hwInfoCpu := cpu.CpuInfo{
		Architecture: "",
		CpuCount:     0,
		Vendor:       vendorId,
		Models: []cpu.Model{
			{},
		},
	}

	result, err := checkCpus(stackDevice, hwInfoCpu)
	if err != nil {
		t.Error(err)
	}
	if result == 0 {
		t.Fatal("CPU vendor should match")
	}

	vendorId = "AuthenticAMD"

	result, err = checkCpus(stackDevice, hwInfoCpu)
	if err != nil {
		t.Error(err)
	}
	if result > 0 {
		t.Fatal("CPU vendor should NOT match")
	}

}
