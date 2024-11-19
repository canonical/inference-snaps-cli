package lspci

func PciDevices() ([]PciDevice, error) {
	hostLsPci, err := hostLsPci()
	if err != nil {
		return nil, err
	}
	devices, err := parseLsPci(hostLsPci)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
