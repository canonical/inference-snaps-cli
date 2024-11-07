package memory

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// Code based on https://github.com/dhamith93/systats/blob/master/memory.go
func GetInfoFromProc() (Info, error) {
	var memoryInfo Info

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return memoryInfo, err
	}

	var memTotal, memFree, memAvailable, memBuffers, memCached, memSh, memReclaimable uint64
	var swapTotal uint64

	meminfoSplit := strings.Split(string(data), "\n")
	for _, line := range meminfoSplit {
		lineArr := strings.Fields(line)
		if len(lineArr) == 0 {
			continue
		}
		if lineArr[0] == "MemTotal:" {
			memTotal, _ = strconv.ParseUint(lineArr[1], 10, 64)
		}
		if lineArr[0] == "MemFree:" {
			memFree, _ = strconv.ParseUint(lineArr[1], 10, 64)
		}
		if lineArr[0] == "MemAvailable:" {
			memAvailable, _ = strconv.ParseUint(lineArr[1], 10, 64)
		}
		if lineArr[0] == "Buffers:" {
			memBuffers, _ = strconv.ParseUint(lineArr[1], 10, 64)
		}
		if lineArr[0] == "Cached:" {
			memCached, _ = strconv.ParseUint(lineArr[1], 10, 64)
		}
		if lineArr[0] == "Shmem:" {
			memSh, _ = strconv.ParseUint(lineArr[1], 10, 64)
		}
		if lineArr[0] == "SReclaimable:" {
			memReclaimable, _ = strconv.ParseUint(lineArr[1], 10, 64)
		}
		if lineArr[0] == "SwapTotal::" {
			swapTotal, _ = strconv.ParseUint(lineArr[1], 10, 64)
		}
	}

	// htop does this https://github.com/hishamhm/htop/blob/8af4d9f453ffa2209e486418811f7652822951c6/linux/LinuxProcessList.c#L832
	memCached = memCached + memReclaimable - memSh

	// Calculated
	var memUsed uint64
	if memTotal > 0 {
		memUsed = memTotal - (memFree + memBuffers + memCached)
	}

	log.Printf("Memory total: %f, used: %f, available: %f", float64(memTotal)/1024/1024, float64(memUsed)/1024/1024, float64(memAvailable)/1024/1024)

	memoryInfo.ramTotal = memTotal
	memoryInfo.swapTotal = swapTotal
	return memoryInfo, nil
}
