package disk

func GetInfo() (*map[string]DirStats, error) {
	var info = make(map[string]DirStats)

	directories := []string{
		"/",
		"/var/lib/snapd/snaps", // https://snapcraft.io/docs/system-snap-directory
	}

	for _, dir := range directories {
		dirInfo, err := GetDirStats(dir)
		if err != nil {
			return nil, err
		}
		info[dir] = dirInfo
	}

	return &info, nil
}
