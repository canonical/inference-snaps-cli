package gpu

func Info() ([]Display, error) {
	var gpus []Display
	lsHw, err := hostLsHw()
	if err != nil {
		return nil, err
	}

	gpus, err = parseLsHw(lsHw)
	if err != nil {
		return nil, err
	}

	return gpus, nil
}
