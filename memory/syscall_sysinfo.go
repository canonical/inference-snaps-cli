package memory

import (
	"golang.org/x/sys/unix"
)

func GetSysInfo() (*unix.Sysinfo_t, error) {
	var sysInfo unix.Sysinfo_t
	err := unix.Sysinfo(&sysInfo)
	if err != nil {
		return nil, err
	}
	return &sysInfo, nil
}
