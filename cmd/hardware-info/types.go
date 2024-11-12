package main

import (
	"hardware-info/cpu"
	"hardware-info/disk"
	"hardware-info/gpu"
	"hardware-info/memory"
)

type HwInfo struct {
	Cpu    *cpu.Info                 `json:"cpu,omitempty"`
	Memory *memory.Info              `json:"memory,omitempty"`
	Disk   *map[string]disk.DirStats `json:"disk,omitempty"`
	Gpu    *[]gpu.Display            `json:"gpu,omitempty"`
}
