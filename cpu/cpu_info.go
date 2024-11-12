package cpu

func GetInfo() (*Info, error) {
	hostLsCpu, err := GetHostLsCpu()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := ParseLsCpu(hostLsCpu)
	if err != nil {
		return nil, err
	}

	return cpuInfo, err
}
