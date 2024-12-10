package types

import (
	"github.com/canonical/hardware-info/hardware_info/cpu"
	"github.com/canonical/hardware-info/hardware_info/disk"
	"github.com/canonical/hardware-info/hardware_info/gpu"
	"github.com/canonical/hardware-info/hardware_info/memory"
)

type HwInfo struct {
	Cpu    *cpu.CpuInfo              `json:"cpu,omitempty"`
	Memory *memory.MemoryInfo        `json:"memory,omitempty"`
	Disk   map[string]*disk.DirStats `json:"disk,omitempty"`
	Gpus   []gpu.Gpu                 `json:"gpu,omitempty"`
}
