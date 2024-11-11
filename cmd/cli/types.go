package main

import (
	"hardware-info/pkg/cpu"
	"hardware-info/pkg/disk"
	"hardware-info/pkg/gpu"
	"hardware-info/pkg/memory"
)

type HwInfo struct {
	Cpu    *cpu.Info            `json:"cpu,omitempty"`
	Memory *memory.Info         `json:"memory,omitempty"`
	Disk   *disk.SystemDirsInfo `json:"disk,omitempty"`
	Gpu    *[]gpu.Display       `json:"gpu,omitempty"`
}
