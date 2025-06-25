package disk

import (
	"fmt"

	"github.com/canonical/stack-utils/pkg/types"
)

var directories = []string{
	"/",
	"/var/lib/snapd/snaps", // https://snapcraft.io/docs/system-snap-directory
}

func Info() (map[string]types.DirStats, error) {
	hostDfData, err := hostDf(directories...)
	if err != nil {
		return nil, err
	}
	return InfoFromData(hostDfData)
}

func InfoFromData(dfData string) (map[string]types.DirStats, error) {
	dirInfos, err := parseDf(dfData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse df: %v", err)
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
