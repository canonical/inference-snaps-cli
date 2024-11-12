package disk

import (
	"golang.org/x/sys/unix"
)

// GetDirStats returns a struct with the total, used, free and available bytes for a given directory.
func GetDirStats(path string) (DirStats, error) {
	dirStats := DirStats{}

	var fs unix.Statfs_t
	err := unix.Statfs(path, &fs)
	if err != nil {
		return dirStats, err
	}

	dirStats.Total = fs.Blocks * uint64(fs.Bsize)
	dirStats.Avail = fs.Bavail * uint64(fs.Bsize)
	dirStats.Free = fs.Bfree * uint64(fs.Bsize)
	dirStats.Used = dirStats.Total - dirStats.Free
	return dirStats, err
}
