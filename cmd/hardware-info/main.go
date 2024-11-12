package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"hardware-info/cpu"
	"hardware-info/disk"
	"hardware-info/gpu"
	"hardware-info/memory"
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
		log.Println("Failed to get memory info:", err)
	}
	hwInfo.Memory = memoryInfo

	cpuInfo, err := cpu.GetInfo()
	if err != nil {
		log.Println("Failed to get CPU info:", err)
	}
	hwInfo.Cpu = cpuInfo

	diskInfo, err := disk.GetInfo()
	if err != nil {
		log.Println("Failed to get disk info:", err)
	}
	hwInfo.Disk = diskInfo

	gpuInfo, err := gpu.GetInfo()
	if err != nil {
		log.Println("Failed to get GPU info:", err)
	}
	hwInfo.Gpu = &gpuInfo

	var jsonString []byte
	if prettyOutput {
		jsonString, err = json.MarshalIndent(hwInfo, "", "  ")
		if err != nil {
			log.Fatal("Failed to marshal to JSON:", err)
		}
	} else {
		jsonString, err = json.Marshal(hwInfo)
		if err != nil {
			log.Fatal("Failed to marshal to JSON:", err)
		}
	}

	if fileOutput != "" {
		err = os.WriteFile(fileOutput, jsonString, 0644)
		if err != nil {
			log.Fatal("Failed to write to file:", err)
		}
	} else {
		fmt.Println(string(jsonString))
	}
}
