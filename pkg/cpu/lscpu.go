package cpu

import (
	"os/exec"
	"strconv"
	"strings"
)

type LsCpuInfo struct {
	Architecture string
	CpuCount     int
	Vendor       string
	Models       []LsCpuModel
}

type LsCpuModel struct {
	Name            string
	Family          *int
	Id              int
	ThreadsPerCore  *int
	Sockets         *int
	CoresPerSocket  *int
	Clusters        *int
	CoresPerCluster *int
	//CpuCount int // sockets * cores-per-socket * clusters * cores-per-cluster * threads-per-core
	MaxFreq  float64
	MinFreq  float64
	BogoMips float64
	Flags    []string
}

func GetLocalLscpu() (string, error) {
	out, err := exec.Command("lscpu").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func ParseLscpu(input string) LsCpuInfo {
	cpuInfo := LsCpuInfo{}
	modelIndex := -1

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		label, value, valid := strings.Cut(line, ": ")
		if !valid {
			continue
		}
		value = strings.TrimSpace(value)

		switch label {
		case "Architecture":
			cpuInfo.Architecture = value
		case "CPU(s)":
			if cpuCount, err := strconv.Atoi(value); err == nil {
				cpuInfo.CpuCount = cpuCount
			}
		case "Vendor ID":
			cpuInfo.Vendor = value

		case "Model name":
			modelIndex++
			cpuInfo.Models = append(cpuInfo.Models, LsCpuModel{Name: value})
		case "CPU family":
			if familyId, err := strconv.Atoi(value); err == nil {
				cpuInfo.Models[modelIndex].Family = &familyId
			}
		case "Model":
			if modelId, err := strconv.Atoi(value); err == nil {
				cpuInfo.Models[modelIndex].Id = modelId
			}
		case "Thread(s) per core":
			if threads, err := strconv.Atoi(value); err == nil {
				cpuInfo.Models[modelIndex].ThreadsPerCore = &threads
			}
		case "Core(s) per socket":
			if cores, err := strconv.Atoi(value); err == nil {
				cpuInfo.Models[modelIndex].CoresPerSocket = &cores
			}
		case "Core(s) per cluster":
			if cores, err := strconv.Atoi(value); err == nil {
				cpuInfo.Models[modelIndex].CoresPerCluster = &cores
			}
		case "Socket(s)":
			if sockets, err := strconv.Atoi(value); err == nil {
				cpuInfo.Models[modelIndex].Sockets = &sockets
			}
		case "Cluster(s)":
			if clusters, err := strconv.Atoi(value); err == nil {
				cpuInfo.Models[modelIndex].Clusters = &clusters
			}
		case "CPU max MHz":
			if maxFreq, err := strconv.ParseFloat(value, 64); err == nil {
				cpuInfo.Models[modelIndex].MaxFreq = maxFreq
			}
		case "CPU min MHz":
			if minFreq, err := strconv.ParseFloat(value, 64); err == nil {
				cpuInfo.Models[modelIndex].MinFreq = minFreq
			}
		case "BogoMIPS":
			if bogoMips, err := strconv.ParseFloat(value, 64); err == nil {
				cpuInfo.Models[modelIndex].BogoMips = bogoMips
			}
		case "Flags":
			flags := strings.Fields(value)
			cpuInfo.Models[modelIndex].Flags = flags
		}

	}

	return cpuInfo
}
