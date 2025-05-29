package hardware_info

import (
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/cpu"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/disk"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/gpu"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/memory"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci_peripherals"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Get(friendlyNames bool) (types.HwInfo, error) {
	var hwInfo types.HwInfo

	memoryInfo, err := memory.Info()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.Memory = memoryInfo

	cpus, err := cpu.Info()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.Cpus = cpus

	diskInfo, err := disk.Info()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.Disk = diskInfo

	gpuInfo, err := gpu.Info(friendlyNames)
	if err != nil {
		return hwInfo, err
	}
	hwInfo.Gpus = gpuInfo

	pciPeripherals, err := pci_peripherals.Info(friendlyNames)
	if err != nil {
		return hwInfo, err
	}
	hwInfo.PciPeripherals = pciPeripherals

	return hwInfo, nil
}
