package main

import (
	"github.com/canonical/hardware-info/cpu"
	"github.com/canonical/hardware-info/disk"
	"github.com/canonical/hardware-info/gpu"
	"github.com/canonical/hardware-info/memory"
)

type HwInfo struct {
	Cpu    *cpu.Info                 `json:"cpu,omitempty"`
	Memory *memory.Info              `json:"memory,omitempty"`
	Disk   *map[string]disk.DirStats `json:"disk,omitempty"`
	Gpu    *[]gpu.Display            `json:"gpu,omitempty"`
}
