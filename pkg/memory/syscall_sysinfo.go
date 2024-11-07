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

	memoryInfo.ramTotal = sysInfo.Totalram
	memoryInfo.swapTotal = sysInfo.Totalswap
	return memoryInfo, nil
}
