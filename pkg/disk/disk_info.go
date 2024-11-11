package disk

func GetInfo() (SystemDirsInfo, error) {
	info := SystemDirsInfo{}

	rootInfo, err := GetDirStats("/")
	info.Root = &rootInfo
	if err != nil {
		return info, err
	}

	// https://snapcraft.io/docs/system-snap-directory
	snapsInfo, err := GetDirStats("/var/lib/snapd/snaps")
	info.Snaps = &snapsInfo
	return info, err
}
