package cpu

import (
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Info() ([]types.CpuInfo, error) {
	hostLsCpu, err := hostLsCpu()
	if err != nil {
		return nil, err
	}

	cpus, err := parseLsCpu(hostLsCpu)
	if err != nil {
		return nil, err
	}

	hostProcCpuInfo, err := procCpuInfo()
	if err != nil {
		return nil, err
	}

	cpus, err = enrichCpus(cpus, hostProcCpuInfo)
	if err != nil {
		return nil, err
	}

	return cpus, err
}

func enrichCpus(cpus []types.CpuInfo, procCpuInfo []ProcCpuInfo) ([]types.CpuInfo, error) {
	for i, cpu := range cpus {

		if cpu.Architecture == "arm64" {
			// lscpu reports ARM CPUs with model ID 1, while cpuinfo has a valid part number.
			// Use /proc/cpuinfo to get the correct value
			// Look up lscpu cpu model in cpuinfo based on bogomips as that looks like the only semi-unique field
			for _, cpuInfo := range procCpuInfo {
				if cpuInfo.BogoMips == cpu.BogoMips {
					cpus[i].ModelId = int(procCpuInfo[0].Part)
				}
			}
		}

	}
	return cpus, nil
}
