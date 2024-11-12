package gpu

func GetInfo() ([]Display, error) {
	var gpus []Display
	lsHw, err := GetHostLsHw()
	if err != nil {
		return nil, err
	}

	gpus, err = ParseLsHw(lsHw)
	if err != nil {
		return nil, err
	}

	return gpus, nil
}
