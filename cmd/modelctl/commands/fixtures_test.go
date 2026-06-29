package commands

import (
	"fmt"

	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/cpu"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
	"github.com/canonical/lscompute/pkg/machine/disk"
	"github.com/canonical/lscompute/pkg/machine/memory"
	lstypes "github.com/canonical/lscompute/pkg/machine/types"
)

func strptr(s string) *string                 { return &s }
func u8ptr(v uint8) *uint8                    { return &v }
func hexptr(v lstypes.HexInt) *lstypes.HexInt { return &v }

// machineInfo returns a small, hand-built MachineInfo fixture for the named
// machine. These replace the old hardware_info.GetFromRawData(test_data)
// helper so the modelctl command tests no longer depend on captured fixtures.
func machineInfo(name string) (*machine.MachineInfo, error) {
	switch name {
	case "dummy-machine":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"fpu", "vme", "de"},
			}},
			Memory: memory.MemoryInfo{TotalRam: 67012501504, TotalSwap: 0},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1006451294208, Avail: 943543738368},
			},
			Devices: []any{pci.Device{
				Bus:                  "pci",
				Slot:                 "0000:00:00.0",
				BusNumber:            0x0,
				DeviceClass:          0x600,
				ProgrammingInterface: u8ptr(0),
				VendorId:             0x8086,
				DeviceId:             0x4637,
				SubvendorId:          hexptr(0x103C),
				SubdeviceId:          hexptr(0x89C6),
				FriendlyNames: pci.FriendlyNames{
					VendorName:    strptr("Intel Corporation"),
					SubvendorName: strptr("Hewlett-Packard Company"),
				},
			}},
		}, nil

	case "i7-1165G7":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"sse4_2", "f16c", "fma", "avx", "avx2", "avx512f"},
			}},
		}, nil

	case "xps13-7390":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"sse4_2", "f16c", "fma", "avx", "avx2"},
			}},
		}, nil

	case "mustang":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
			}},
		}, nil

	default:
		return nil, fmt.Errorf("no machine fixture for %q", name)
	}
}
