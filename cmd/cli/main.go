package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"hardware-info/pkg/cpu"
	"hardware-info/pkg/disk"
	"hardware-info/pkg/gpu"
	"hardware-info/pkg/memory"
)

func main() {
	var prettyOutput bool
	var fileOutput string

	flag.BoolVar(&prettyOutput, "pretty", false, "Output pretty json. Default is compact json.")
	flag.StringVar(&fileOutput, "file", "", "Output json to this file. Default output is to stdout.")
	flag.Parse()

	hwInfo := HwInfo{}

	memoryInfo, err := memory.GetInfo()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get memory info: %s\n", err)
	}
	hwInfo.Memory = &memoryInfo

	cpuInfo, err := cpu.GetInfo()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get cpu info: %s\n", err)
	}
	hwInfo.Cpu = &cpuInfo

	diskInfo, err := disk.GetInfo()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get disk info: %s\n", err)
	}
	hwInfo.Disk = &diskInfo

	gpuInfo, err := gpu.GetInfo()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get gpu info: %s\n", err)
	}
	hwInfo.Gpu = &gpuInfo

	var jsonString []byte
	if prettyOutput {
		jsonString, err = json.MarshalIndent(hwInfo, "", "  ")
	} else {
		jsonString, _ = json.Marshal(hwInfo)
	}

	if fileOutput != "" {
		err = os.WriteFile(fileOutput, jsonString, 0644)
	} else {
		fmt.Println(string(jsonString))
	}
}
