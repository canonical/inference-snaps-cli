package selector

import (
	"testing"

	types2 "github.com/canonical/ml-snap-utils/pkg/hardware_info/types"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func TestCheckCpu(t *testing.T) {
	vendorId := "GenuineIntel"
	stackDevice := types.StackDevice{
		Type:     "cpu",
		Bus:      nil,
		VendorId: &vendorId,
	}

	hwInfoCpu := types2.CpuInfo{
		Architecture: "",
		CpuCount:     0,
		Vendor:       vendorId,
		Models: []types2.Model{
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
