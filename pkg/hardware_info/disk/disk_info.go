package disk

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

// We assume this is always running inside a snap, so we look at the snap data directory's available space
const directory = "/var/lib/snapd/snaps" // https://snapcraft.io/docs/system-snap-directory

// Info returns the total size and available size using the statfs syscall.
func Info() (*types.DirStats, error) {
	dirInfo, err := statFs(directory)
	if err != nil {
		return nil, fmt.Errorf("error getting directory info: %v", err)
	}

	return &dirInfo, nil
}

// InfoFromRawData parses the output of the df command and returns the total size and available size
func InfoFromRawData(dfData string) (*types.DirStats, error) {
	dirInfos, err := parseDf(dfData)
	if err != nil {
		return nil, fmt.Errorf("error parsing df: %v", err)
	}

	return &dirInfos[0], nil
}
