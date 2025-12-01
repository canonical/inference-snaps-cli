package disk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

var directories = []string{
	"/",
	"/var/lib/snapd/snaps", // https://snapcraft.io/docs/system-snap-directory
}

// Info returns the total size and available size for root and snap dirs on the host system, using the statfs syscall.
func Info() (map[string]types.DirStats, error) {
	var info = make(map[string]types.DirStats)

	for _, dir := range directories {
		dirInfo, err := statFs(dir)
		if err != nil {
			return nil, fmt.Errorf("error getting directory info: %v", err)
		}
		info[dir] = dirInfo
	}

	return info, nil
}

// InfoFromRawData returns the total size and available size of the root and snap dirs, taking a string in which represents
// the  output of the df command.
func InfoFromRawData(dfData string) (map[string]types.DirStats, error) {
	dirInfos, err := parseDf(dfData)
	if err != nil {
		return nil, fmt.Errorf("error parsing df: %v", err)
	}

	if len(dirInfos) != len(directories) {
		return nil, fmt.Errorf("df did not return info for all dirs")
	}

	var info = make(map[string]types.DirStats)
	for i, dir := range directories {
		info[dir] = dirInfos[i]
	}

	return info, nil
}

func parseDf(dfData string) ([]types.DirStats, error) {
	var parsedDirStats []types.DirStats

	lines := strings.Split(dfData, "\n")

	// Skip header line
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) != 6 {
			return nil, fmt.Errorf("not 6 columns")
		}

		totalSize, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'total blocks' field: %v", err)
		}
		//usedSize, err := strconv.ParseUint(fields[2], 10, 64)
		availableSize, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'available blocks' field: %v", err)
		}

		var thisDir = types.DirStats{
			Total: totalSize,
			Avail: availableSize,
		}
		parsedDirStats = append(parsedDirStats, thisDir)
	}

	return parsedDirStats, nil
}
