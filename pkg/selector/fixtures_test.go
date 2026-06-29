package selector

import (
	"fmt"

	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/cpu"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
	"github.com/canonical/lscompute/pkg/machine/disk"
	"github.com/canonical/lscompute/pkg/machine/memory"
	lstypes "github.com/canonical/lscompute/pkg/machine/types"
)

func fxStr(s string) *string                 { return &s }
func fxU8(v uint8) *uint8                    { return &v }
func fxHex(v lstypes.HexInt) *lstypes.HexInt { return &v }

// machineInfo returns an in-code MachineInfo fixture for the named machine.
// PCI devices are trimmed to those relevant to engine matching (GPUs,
// accelerators, or devices carrying additional properties).
func machineInfo(name string) (*machine.MachineInfo, error) {
	switch name {
	case "ampere-one-m-banshee-12":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0xC0,
					PartNumber:    0xAC4,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "jscvt", "fcma", "lrcpc", "dcpop", "sha3", "sm3", "sm4", "asimddp", "sha512", "asimdfhm", "dit", "uscat", "ilrcpc", "flagm", "ssbs", "sb", "paca", "pacg", "dcpodp", "flagm2", "frint", "i8mm", "bf16", "rng", "bti", "ecv"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32876765184, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 502392942592, Avail: 431927586816},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0004:02:00.0",
					BusNumber:            0x2,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x1A03,
					DeviceId:             0x2000,
					SubvendorId:          fxHex(0x1A03),
					SubdeviceId:          fxHex(0x2000),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("ASPEED Technology, Inc."),
						DeviceName:    fxStr("ASPEED Graphics Family"),
						SubvendorName: fxStr("ASPEED Technology, Inc."),
						SubdeviceName: fxStr("Onboard Graphics"),
					},
				},
			},
		}, nil
	case "ampere-one-siryn":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0xC0,
					PartNumber:    0xAC3,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "jscvt", "fcma", "lrcpc", "dcpop", "sha3", "asimddp", "sha512", "asimdfhm", "dit", "uscat", "ilrcpc", "flagm", "ssbs", "sb", "paca", "pacg", "dcpodp", "flagm2", "frint", "i8mm", "bf16", "rng", "bti", "ecv"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32876765184, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 502392942592, Avail: 431927586816},
			},
		}, nil
	case "ampere-one-x-banshee-8":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0xC0,
					PartNumber:    0xAC4,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "jscvt", "fcma", "lrcpc", "dcpop", "sha3", "sm3", "sm4", "asimddp", "sha512", "asimdfhm", "dit", "uscat", "ilrcpc", "flagm", "ssbs", "sb", "paca", "pacg", "dcpodp", "flagm2", "frint", "i8mm", "bf16", "rng", "bti", "ecv"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32876765184, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 502392942592, Avail: 431927586816},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0004:03:00.0",
					BusNumber:            0x3,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x1A03,
					DeviceId:             0x2000,
					SubvendorId:          fxHex(0x1A03),
					SubdeviceId:          fxHex(0x2000),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("ASPEED Technology, Inc."),
						DeviceName:    fxStr("ASPEED Graphics Family"),
						SubvendorName: fxStr("ASPEED Technology, Inc."),
						SubdeviceName: fxStr("Onboard Graphics"),
					},
				},
			},
		}, nil
	case "asus-ux301l":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "cpuid_fault", "epb", "pti", "ssbd", "ibrs", "ibpb", "stibp", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "xsaveopt", "dtherm", "ida", "arat", "pln", "pts", "vnmi", "md_clear", "flush_l1d"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 7663951872, TotalSwap: 4294963200},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 124313583616, Avail: 61499084800},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0xA2E,
					SubvendorId:          fxHex(0x1043),
					SubdeviceId:          fxHex(0x13BD),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Haswell-ULT Integrated Graphics Controller"),
						SubvendorName: fxStr("ASUSTeK Computer Inc."),
					},
				},
			},
		}, nil
	case "hp-pavilion-15-cs-3037nl":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "avx512f", "avx512dq", "rdseed", "adx", "smap", "avx512ifma", "clflushopt", "intel_pt", "avx512cd", "sha_ni", "avx512bw", "avx512vl", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "vnmi", "avx512vbmi", "umip", "pku", "ospke", "avx512_vbmi2", "gfni", "vaes", "vpclmulqdq", "avx512_vnni", "avx512_bitalg", "avx512_vpopcntdq", "rdpid", "fsrm", "md_clear", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 16537874432, TotalSwap: 19327348736},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 498347700224, Avail: 61093965824},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x8A52,
					SubvendorId:          fxHex(0x103C),
					SubdeviceId:          fxHex(0x86E2),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Iris Plus Graphics G7"),
						SubvendorName: fxStr("Hewlett-Packard Company"),
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:06:00.0",
					BusNumber:            0x6,
					DeviceClass:          0x302,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x10DE,
					DeviceId:             0x1D13,
					SubvendorId:          fxHex(0x103C),
					SubdeviceId:          fxHex(0x86E2),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("NVIDIA Corporation"),
						DeviceName:    fxStr("GP108M [GeForce MX250]"),
						SubvendorName: fxStr("Hewlett-Packard Company"),
					},
					AdditionalProperties: map[string]string{
						"computeCapability": "6.1",
						"vram":              "2147483648",
					},
				},
			},
		}, nil
	case "hp-proliant-rl300-gen11-altra":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0x41,
					PartNumber:    0xD0C,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "lrcpc", "dcpop", "asimddp", "ssbs"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 33134252032, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 470888099840, Avail: 435367202816},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0006:01:00.1",
					BusNumber:            0x1,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x102B,
					DeviceId:             0x538,
					SubvendorId:          fxHex(0x1590),
					SubdeviceId:          fxHex(0xE4),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Matrox Electronics Systems Ltd."),
						DeviceName:    fxStr("MGA G200eH3"),
						SubvendorName: fxStr("Hewlett Packard Enterprise"),
						SubdeviceName: fxStr("iLO5 VGA"),
					},
				},
			},
		}, nil
	case "hp-proliant-rl300-gen11-altra-max":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0x41,
					PartNumber:    0xD0C,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "lrcpc", "dcpop", "asimddp", "ssbs"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 66839527424, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 470888099840, Avail: 425825845248},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0006:01:00.1",
					BusNumber:            0x1,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x102B,
					DeviceId:             0x538,
					SubvendorId:          fxHex(0x1590),
					SubdeviceId:          fxHex(0xE4),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Matrox Electronics Systems Ltd."),
						DeviceName:    fxStr("MGA G200eH3"),
						SubvendorName: fxStr("Hewlett Packard Enterprise"),
						SubdeviceName: fxStr("iLO5 VGA"),
					},
				},
			},
		}, nil
	case "hp-zbook-power-16-inch-g11":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "intel_ppin", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 33047601152, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1003736440832, Avail: 876099534848},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x7DD5,
					SubvendorId:          fxHex(0x103C),
					SubdeviceId:          fxHex(0x8CA0),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Meteor Lake-P [Intel Graphics]"),
						SubvendorName: fxStr("Hewlett-Packard Company"),
					},
					AdditionalProperties: map[string]string{
						"vram": "30592884736",
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:0b.0",
					BusNumber:            0x0,
					DeviceClass:          0x1200,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x7D1D,
					SubvendorId:          fxHex(0x103C),
					SubdeviceId:          fxHex(0x8CA0),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Meteor Lake NPU"),
						SubvendorName: fxStr("Hewlett-Packard Company"),
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:01:00.0",
					BusNumber:            0x1,
					DeviceClass:          0x302,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x10DE,
					DeviceId:             0x28B9,
					SubvendorId:          fxHex(0x103C),
					SubdeviceId:          fxHex(0x8CA3),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("NVIDIA Corporation"),
						DeviceName:    fxStr("AD107GLM [RTX 1000 Ada Generation Laptop GPU]"),
						SubvendorName: fxStr("Hewlett-Packard Company"),
					},
				},
			},
		}, nil
	case "hp-zbook-i712850HX+RadeonPROW6600M":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "tme", "rdpid", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities", "ibpb_exit_to_user"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 67012501504, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1006451294208, Avail: 943543738368},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x4688,
					SubvendorId:          fxHex(0x103C),
					SubdeviceId:          fxHex(0x89C6),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Alder Lake-HX GT1 [UHD Graphics 770]"),
						SubvendorName: fxStr("Hewlett-Packard Company"),
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:03:00.0",
					BusNumber:            0x3,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x1002,
					DeviceId:             0x73E1,
					SubvendorId:          fxHex(0x103C),
					SubdeviceId:          fxHex(0x89C6),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Advanced Micro Devices, Inc. [AMD/ATI]"),
						DeviceName:    fxStr("Navi 23 WKS-XM [Radeon PRO W6600M]"),
						SubvendorName: fxStr("Hewlett-Packard Company"),
					},
					AdditionalProperties: map[string]string{
						"microarchitecture": "gfx1032",
						"vram":              "8573157376",
					},
				},
			},
		}, nil
	case "i5-3570k+arc-a580+gtx1080ti":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "rdtscp", "lm", "constant_tsc", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "est", "tm2", "ssse3", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "cpuid_fault", "epb", "pti", "ssbd", "ibrs", "ibpb", "stibp", "fsgsbase", "smep", "erms", "xsaveopt", "dtherm", "ida", "arat", "pln", "pts", "md_clear", "flush_l1d"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 16652423168, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1005867986944, Avail: 862075604992},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x162,
					SubvendorId:          fxHex(0x1458),
					SubdeviceId:          fxHex(0xD000),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("IvyBridge GT2 [HD Graphics 4000]"),
						SubvendorName: fxStr("Gigabyte Technology Co., Ltd"),
					},
					AdditionalProperties: map[string]string{
						"vram": "17179869184",
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:01:00.0",
					BusNumber:            0x1,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x10DE,
					DeviceId:             0x1B06,
					SubvendorId:          fxHex(0x10DE),
					SubdeviceId:          fxHex(0x1B06),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("NVIDIA Corporation"),
						DeviceName:    fxStr("GP102 [GeForce GTX 1080 Ti]"),
						SubvendorName: fxStr("NVIDIA Corporation"),
					},
					AdditionalProperties: map[string]string{
						"compute_capability": "6.1",
						"vram":               "11811160064",
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:04:00.0",
					BusNumber:            0x4,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x56A2,
					SubvendorId:          fxHex(0x1849),
					SubdeviceId:          fxHex(0x6003),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("DG2 [Arc A580]"),
						SubvendorName: fxStr("ASRock Incorporation"),
					},
					AdditionalProperties: map[string]string{
						"vram": "8096681984",
					},
				},
			},
		}, nil
	case "i7-1165G7":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "cat_l2", "invpcid_single", "cdp_l2", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "vnmi", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdt_a", "avx512f", "avx512dq", "rdseed", "adx", "smap", "avx512ifma", "clflushopt", "clwb", "intel_pt", "avx512cd", "sha_ni", "avx512bw", "avx512vl", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "avx512vbmi", "umip", "pku", "ospke", "avx512_vbmi2", "gfni", "vaes", "vpclmulqdq", "avx512_vnni", "avx512_bitalg", "avx512_vpopcntdq", "rdpid", "movdiri", "movdir64b", "fsrm", "avx512_vp2intersect", "md_clear", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 16493559808, TotalSwap: 1027600384},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 500609097728, Avail: 20873912320},
			},
			Devices: []any{
				pci.Device{
					Bus:         "pci",
					Slot:        "0000:00:02.0",
					BusNumber:   0x0,
					DeviceClass: 0x300,
					VendorId:    0x8086,
					DeviceId:    0x9A49,
					SubvendorId: fxHex(0x1028),
					SubdeviceId: fxHex(0xB20),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("TigerLake-LP GT2 [Iris Xe Graphics]"),
						SubvendorName: fxStr("Dell"),
					},
				},
			},
		}, nil
	case "i7-2600k+arc-a580":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ht", "tm", "pbe", "syscall", "nx", "rdtscp", "lm", "constant_tsc", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "lahf_lm", "epb", "pti", "ssbd", "ibrs", "ibpb", "stibp", "tpr_shadow", "flexpriority", "ept", "vpid", "xsaveopt", "dtherm", "ida", "arat", "pln", "pts", "vnmi", "md_clear", "flush_l1d"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 11439935488, TotalSwap: 2147479552},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 58225459200, Avail: 16706797568},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:03:00.0",
					BusNumber:            0x3,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x56A2,
					SubvendorId:          fxHex(0x1849),
					SubdeviceId:          fxHex(0x6003),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("DG2 [Arc A580]"),
						SubvendorName: fxStr("ASRock Incorporation"),
					},
					AdditionalProperties: map[string]string{
						"vram": "8096681984",
					},
				},
			},
		}, nil
	case "i7-10510U":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "sgx", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "mpx", "rdseed", "adx", "smap", "clflushopt", "intel_pt", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "vnmi", "md_clear", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 15908761600, TotalSwap: 4294963200},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1003736440832, Avail: 732326371328},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x9B41,
					SubvendorId:          fxHex(0x1028),
					SubdeviceId:          fxHex(0x962),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("CometLake-U GT2 [UHD Graphics]"),
						SubvendorName: fxStr("Dell"),
					},
				},
			},
		}, nil
	case "mustang":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "tme", "rdpid", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32876765184, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 502392942592, Avail: 431927586816},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:03:00.0",
					BusNumber:            0x3,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0xE20B,
					SubvendorId:          fxHex(0x8086),
					SubdeviceId:          fxHex(0x1100),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Battlemage G21 [Arc B580]"),
						SubvendorName: fxStr("Intel Corporation"),
					},
					AdditionalProperties: map[string]string{
						"vram": "12168933376",
					},
				},
			},
		}, nil
	case "raspberry-pi-5":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0x41,
					PartNumber:    0xD0B,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "lrcpc", "dcpop", "asimddp"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 8323276800, TotalSwap: 0},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 62270910464, Avail: 40666140672},
			},
		}, nil
	case "raspberry-pi-5+hailo-8":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0x41,
					PartNumber:    0xD0B,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "lrcpc", "dcpop", "asimddp"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 8317198336, TotalSwap: 1073737728},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 30777720832, Avail: 16953257984},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0001:01:00.0",
					BusNumber:            0x1,
					DeviceClass:          0xB40,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x1E60,
					DeviceId:             0x2864,
					SubvendorId:          fxHex(0x1E60),
					SubdeviceId:          fxHex(0x2864),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Hailo Technologies Ltd."),
						DeviceName:    fxStr("Hailo-8 AI Processor"),
						SubvendorName: fxStr("Hailo Technologies Ltd."),
					},
				},
			},
		}, nil
	case "system76-addw4":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32945307648, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 980105256960, Avail: 574739996672},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0xA788,
					SubvendorId:          fxHex(0x1558),
					SubdeviceId:          fxHex(0x353),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Raptor Lake-S UHD Graphics"),
						SubvendorName: fxStr("CLEVO/KAPOK Computer"),
					},
					AdditionalProperties: map[string]string{
						"vram": "31002615808",
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:01:00.0",
					BusNumber:            0x1,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x10DE,
					DeviceId:             0x28A1,
					SubvendorId:          fxHex(0x1558),
					SubdeviceId:          fxHex(0x353),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("NVIDIA Corporation"),
						DeviceName:    fxStr("AD107M [GeForce RTX 4050 Max-Q / Mobile]"),
						SubvendorName: fxStr("CLEVO/KAPOK Computer"),
					},
					AdditionalProperties: map[string]string{
						"compute_capability": "8.9",
						"vram":               "6439305216",
					},
				},
			},
		}, nil
	case "xps13-7390":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "sgx", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "mpx", "rdseed", "adx", "smap", "clflushopt", "intel_pt", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "vnmi", "md_clear", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 15908761600, TotalSwap: 4294963200},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1003736440832, Avail: 732343107584},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x9B41,
					SubvendorId:          fxHex(0x1028),
					SubdeviceId:          fxHex(0x962),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("CometLake-U GT2 [UHD Graphics]"),
						SubvendorName: fxStr("Dell"),
					},
					AdditionalProperties: map[string]string{
						"vram": "14482374656",
					},
				},
			},
		}, nil
	case "xps13-9350":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdt_a", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "lam", "wbnoinvd", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "intel_ppin", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdt_a", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "lam", "wbnoinvd", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 33196507136, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 993185353728, Avail: 922389172224},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x64A0,
					SubvendorId:          fxHex(0x1028),
					SubdeviceId:          fxHex(0xCC9),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Lunar Lake [Intel Arc Graphics 130V / 140V]"),
						SubvendorName: fxStr("Dell"),
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:0b.0",
					BusNumber:            0x0,
					DeviceClass:          0x1200,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x8086,
					DeviceId:             0x643E,
					SubvendorId:          fxHex(0x1028),
					SubdeviceId:          fxHex(0xCC9),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Intel Corporation"),
						DeviceName:    fxStr("Lunar Lake NPU"),
						SubvendorName: fxStr("Dell"),
					},
				},
			},
		}, nil
	case "lenovo-thinkpad-p16s":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "AuthenticAMD",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "mmx", "fxsr", "sse", "sse2", "ht", "syscall", "nx", "mmxext", "fxsr_opt", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "rep_good", "amd_lbr_v2", "nopl", "xtopology", "nonstop_tsc", "cpuid", "extd_apicid", "aperfmperf", "rapl", "pni", "pclmulqdq", "monitor", "ssse3", "fma", "cx16", "sse4_1", "sse4_2", "movbe", "popcnt", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "cmp_legacy", "svm", "extapic", "cr8_legacy", "abm", "sse4a", "misalignsse", "3dnowprefetch", "osvw", "ibs", "skinit", "wdt", "tce", "topoext", "perfctr_core", "perfctr_nb", "bpext", "perfctr_llc", "mwaitx", "cpuid_fault", "cpb", "cat_l3", "cdp_l3", "hw_pstate", "ssbd", "mba", "perfmon_v2", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "vmmcall", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "cqm", "rdt_a", "avx512f", "avx512dq", "adx", "smap", "avx512ifma", "clflushopt", "clwb", "avx512cd", "sha_ni", "avx512bw", "avx512vl", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "cqm_llc", "cqm_occup_llc", "cqm_mbm_total", "cqm_mbm_local", "user_shstk", "avx_vnni", "avx512_bf16", "clzero", "irperf", "xsaveerptr", "rdpru", "wbnoinvd", "cppc", "arat", "npt", "lbrv", "svm_lock", "nrip_save", "tsc_scale", "vmcb_clean", "flushbyasid", "decodeassists", "pausefilter", "pfthreshold", "avic", "v_vmsave_vmload", "vgif", "x2avic", "v_spec_ctrl", "vnmi", "avx512vbmi", "umip", "pku", "ospke", "avx512_vbmi2", "gfni", "vaes", "vpclmulqdq", "avx512_vnni", "avx512_bitalg", "avx512_vpopcntdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "overflow_recov", "succor", "smca", "fsrm", "avx512_vp2intersect", "flush_l1d", "amd_lbr_pmc_freeze"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 92413652992, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 2011859927040, Avail: 521356337152},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:c4:00.0",
					BusNumber:            0xC4,
					DeviceClass:          0x380,
					ProgrammingInterface: fxU8(0),
					VendorId:             0x1002,
					DeviceId:             0x1114,
					SubvendorId:          fxHex(0x17AA),
					SubdeviceId:          fxHex(0x512F),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    fxStr("Advanced Micro Devices, Inc. [AMD/ATI]"),
						DeviceName:    fxStr("Krackan [Radeon 840M / 860M Graphics]"),
						SubvendorName: fxStr("Lenovo"),
					},
					AdditionalProperties: map[string]string{
						"microarchitecture": "gfx1152",
						"vram":              "8589934592",
					},
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("no machine fixture for %q", name)
	}
}
