package memory

func GetInfo() (*Info, error) {
	var memoryInfo Info

	sysInfo, err := GetSysInfo()
	if err != nil {
		return nil, err
	}

	// The memory size fields need to be multiplied by the unit to get to bytes
	memoryInfo.RamTotal = sysInfo.Totalram * uint64(sysInfo.Unit)
	memoryInfo.SwapTotal = sysInfo.Totalswap * uint64(sysInfo.Unit)
	return &memoryInfo, nil
}
