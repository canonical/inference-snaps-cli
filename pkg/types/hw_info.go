package types

import (
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/types"
)

type HwInfo struct {
	Cpus   []types.CpuInfo            `json:"cpu,omitempty"`
	Memory *types.MemoryInfo          `json:"memory,omitempty"`
	Disk   map[string]*types.DirStats `json:"disk,omitempty"`
	Gpus   []types.Gpu                `json:"gpu,omitempty"`
}
