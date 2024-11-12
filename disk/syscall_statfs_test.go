package disk

import (
	"testing"

	"hardware-info/utils"
)

var testDirs = []string{
	"/",
	"/var/lib/snapd/snaps",
}

func TestGetDirStats(t *testing.T) {
	for _, dir := range testDirs {
		t.Run(dir, func(t *testing.T) {
			diskStats, err := GetDirStats(dir)
			if err != nil {
				t.Fatalf("GetDirStats() failed: %v", err)
			}

			t.Log("Total:", utils.FmtGigabytes(diskStats.Total))
			t.Log("Used:", utils.FmtGigabytes(diskStats.Used))
			t.Log("Avail:", utils.FmtGigabytes(diskStats.Avail))
		})
	}
}
