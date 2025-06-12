package cpu

import (
	"fmt"
	"slices"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Info() ([]types.CpuInfo, error) {

	hostProcCpuInfo, err := procCpuInfo()
	if err != nil {
		return nil, err
	}

	cpus, err := copyCpuInfoToStruct(hostProcCpuInfo)

	return cpus, nil
}

func copyCpuInfoToStruct(procCpus []ProcCpuInfo) ([]types.CpuInfo, error) {
	procCpus = slices.CompactFunc(procCpus, procCpuInfoIsDuplicate)

	cpuInfos, err := cpuInfoFromProc(procCpus)
	if err != nil {
		return nil, err
	}
	return cpuInfos, nil
}

func procCpuInfoIsDuplicate(a ProcCpuInfo, b ProcCpuInfo) bool {
	if a.Architecture != b.Architecture {
		return false
	}

	if a.Architecture == amd64 {
		if a.ManufacturerId != b.ManufacturerId {
			return false
		}
		if a.BrandString != b.BrandString {
			return false
		}
		for _, flag := range a.Flags {
			if !slices.Contains(b.Flags, flag) {
				return false
			}
		}
		return true
	}

	if a.Architecture == arm64 {
		if a.ImplementerId != b.ImplementerId {
			return false
		}
		if a.PartNumber != b.PartNumber {
			return false
		}
		if a.Variant != b.Variant {
			return false
		}
		if a.Revision != b.Revision {
			return false
		}

		for _, feature := range a.Features {
			if !slices.Contains(b.Features, feature) {
				return false
			}
		}

		return true
	}
	return false
}

func cpuInfoFromProc(procCpus []ProcCpuInfo) ([]types.CpuInfo, error) {
	var cpuInfos []types.CpuInfo
	for _, procCpu := range procCpus {
		var cpuInfo types.CpuInfo
		if procCpu.Architecture == amd64 {
			cpuInfo.Architecture = procCpu.Architecture
			cpuInfo.ManufacturerId = procCpu.ManufacturerId
			cpuInfo.Flags = procCpu.Flags
		} else if procCpu.Architecture == arm64 {
			cpuInfo.Architecture = procCpu.Architecture
			cpuInfo.ImplementerId = types.HexInt(procCpu.ImplementerId)
			cpuInfo.PartNumber = types.HexInt(procCpu.PartNumber)
		} else {
			return nil, fmt.Errorf("unsupported architecture %s", procCpu.Architecture)
		}
		cpuInfos = append(cpuInfos, cpuInfo)
	}
	return cpuInfos, nil
}
