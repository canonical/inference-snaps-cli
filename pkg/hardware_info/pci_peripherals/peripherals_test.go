package pci_peripherals

import (
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
)

func TestAccelerators(t *testing.T) {
	pciDevices := []pci.PciDevice{
		{
			DeviceClass: 0x3012, // mock GPU
		},
	}
	pciAccelerators, err := accelerators(pciDevices)
	if err != nil {
		t.Fatal(err)
	}
	if len(pciAccelerators) != 0 {
		t.Errorf("Accelerators returned invalid device")
	}

	pciDevices = append(pciDevices, pci.PciDevice{DeviceClass: 0x1200}) // mock NPU
	pciAccelerators, err = accelerators(pciDevices)
	if err != nil {
		t.Fatal(err)
	}
	if len(pciAccelerators) != 1 {
		t.Errorf("Accelerators did not return valid device")
	}
}

func TestCoprocessors(t *testing.T) {
	pciDevices := []pci.PciDevice{
		{
			DeviceClass: 0x3000, // mock GPU
		},
	}
	pciAccelerators, err := coprocessors(pciDevices)
	if err != nil {
		t.Fatal(err)
	}
	if len(pciAccelerators) != 0 {
		t.Errorf("Coprocessors returned invalid device")
	}

	pciDevices = append(pciDevices, pci.PciDevice{DeviceClass: 0x0b40}) // mock Hailo accelerator
	pciAccelerators, err = coprocessors(pciDevices)
	if err != nil {
		t.Fatal(err)
	}
	if len(pciAccelerators) != 1 {
		t.Errorf("Coprocessors did not return valid device")
	}
}
