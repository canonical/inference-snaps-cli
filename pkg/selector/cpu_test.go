package selector

import (
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func TestCheckCpuVendor(t *testing.T) {
	vendorId := "GenuineIntel"
	stackDevice := types.StackDevice{
		Type:     "cpu",
		Bus:      nil,
		VendorId: &vendorId,
	}

	hwInfoCpus := []types.CpuInfo{{
		Architecture: "",
		VendorId:     vendorId,
	}}

	compatible, reason, err := checkCpus(stackDevice, hwInfoCpus)
	if err != nil {
		t.Fatal(err)
	}
	if !compatible {
		t.Fatalf("CPU vendor should match: %s", reason)
	}

	vendorId = "AuthenticAMD"

	compatible, reason, err = checkCpus(stackDevice, hwInfoCpus)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatal("CPU vendor should NOT match")
	}

}

func TestCheckCpuFlags(t *testing.T) {
	vendorId := "GenuineIntel"
	stackDevice := types.StackDevice{
		Type:     "cpu",
		Bus:      nil,
		VendorId: &vendorId,
		Flags:    []string{"avx2"},
	}

	hwInfoCpus := []types.CpuInfo{{
		Architecture: "",
		VendorId:     vendorId,
		Flags:        []string{"avx2"},
	}}

	compatible, reason, err := checkCpus(stackDevice, hwInfoCpus)
	if err != nil {
		t.Fatal(err)
	}
	if !compatible {
		t.Fatalf("CPU flags should match: %s", reason)
	}

	stackDevice.Flags = []string{"avx512"}

	compatible, reason, err = checkCpus(stackDevice, hwInfoCpus)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatal("CPU flags should NOT match")
	}

}
