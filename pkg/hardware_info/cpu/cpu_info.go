package cpu

import (
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/types"
)

func Info() (*types.CpuInfo, error) {
	hostLsCpu, err := hostLsCpu()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := parseLsCpu(hostLsCpu)
	if err != nil {
		return nil, err
	}

	return cpuInfo, err
}
