package memory

import (
	"golang.org/x/sys/unix"
)

func GetInfo() (Info, error) {
	var memoryInfo Info

	var sysInfo unix.Sysinfo_t
	err := unix.Sysinfo(&sysInfo)
	if err != nil {
		return memoryInfo, err
	}

	// The memory size fields need to be multiplied by the unit to get to bytes
	memoryInfo.RamTotal = sysInfo.Totalram * uint64(sysInfo.Unit)
	memoryInfo.SwapTotal = sysInfo.Totalswap * uint64(sysInfo.Unit)
	return memoryInfo, nil
}
