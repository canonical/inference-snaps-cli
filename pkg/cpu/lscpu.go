package cpu

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

type LscpuContainer struct {
	Lscpu []LscpuObject `json:"lscpu"`
}

type LscpuObject struct {
	Field    string        `json:"field"`
	Data     string        `json:"data"`
	Children []LscpuObject `json:"children"`
}

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

func GetInfo() (LsCpuInfo, error) {
	localLscpu, err := GetLocalLscpu()
	if err != nil {
		return LsCpuInfo{}, err
	}

	return ParseLscpu(localLscpu), nil

}

func GetLocalLscpu() ([]byte, error) {
	out, err := exec.Command("lscpu", "--json", "--hierarchic").Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ParseLscpu(input []byte) LsCpuInfo {
	cpuInfo := LsCpuInfo{}

	var lscpuJson LscpuContainer
	err := json.Unmarshal(input, &lscpuJson)
	if err != nil {
		return LsCpuInfo{}
	}

	for _, line := range lscpuJson.Lscpu {
		label := line.Field
		value := line.Data

		switch label {
		case "Architecture:":
			cpuInfo.Architecture = value
		case "CPU(s):":
			if cpuCount, err := strconv.Atoi(value); err == nil {
				cpuInfo.CpuCount = cpuCount
			}
		case "Vendor ID:":
			cpuInfo.Vendor = value

			for _, vendorChild := range line.Children {
				switch vendorChild.Field {

				case "Model name:":
					modelOutput := LsCpuModel{Name: value}
					modelOutput.Name = vendorChild.Data

					for _, modelNameChild := range vendorChild.Children {
						switch modelNameChild.Field {
						case "CPU family:":
							if familyId, err := strconv.Atoi(modelNameChild.Data); err == nil {
								modelOutput.Family = &familyId
							}
						case "Model:":
							if modelId, err := strconv.Atoi(modelNameChild.Data); err == nil {
								modelOutput.Id = modelId
							}
						case "Thread(s) per core:":
							if threads, err := strconv.Atoi(modelNameChild.Data); err == nil {
								modelOutput.ThreadsPerCore = &threads
							}
						case "Core(s) per socket:":
							if cores, err := strconv.Atoi(modelNameChild.Data); err == nil {
								modelOutput.CoresPerSocket = &cores
							}
						case "Core(s) per cluster:":
							if cores, err := strconv.Atoi(modelNameChild.Data); err == nil {
								modelOutput.CoresPerCluster = &cores
							}
						case "Socket(s):":
							if sockets, err := strconv.Atoi(modelNameChild.Data); err == nil {
								modelOutput.Sockets = &sockets
							}
						case "Cluster(s):":
							if clusters, err := strconv.Atoi(modelNameChild.Data); err == nil {
								modelOutput.Clusters = &clusters
							}
						case "CPU max MHz:":
							if maxFreq, err := strconv.ParseFloat(modelNameChild.Data, 64); err == nil {
								modelOutput.MaxFreq = maxFreq
							}
						case "CPU min MHz:":
							if minFreq, err := strconv.ParseFloat(modelNameChild.Data, 64); err == nil {
								modelOutput.MinFreq = minFreq
							}
						case "BogoMIPS:":
							if bogoMips, err := strconv.ParseFloat(modelNameChild.Data, 64); err == nil {
								modelOutput.BogoMips = bogoMips
							}
						case "Flags:":
							flags := strings.Fields(modelNameChild.Data)
							modelOutput.Flags = flags
						}
					}
					cpuInfo.Models = append(cpuInfo.Models, modelOutput)
				}
			}
		}
	}

	return cpuInfo
}
