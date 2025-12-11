package main

import (
	"fmt"
	"log"
	"os"

	"github.com/canonical/inference-snaps-cli/pkg/hardware_info/pci/intel"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: test-vram <pci slot address>")
	}

	firstArg := os.Args[1]

	device := types.PciDevice{
		Slot:                 firstArg,
		BusNumber:            0,
		DeviceClass:          0,
		ProgrammingInterface: nil,
		VendorId:             0,
		DeviceId:             0,
		SubvendorId:          nil,
		SubdeviceId:          nil,
		PciFriendlyNames:     types.PciFriendlyNames{},
		AdditionalProperties: nil,
	}

	vram, err := intel.VRamOpenCl(device)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("vRam: ", *vram)
}
