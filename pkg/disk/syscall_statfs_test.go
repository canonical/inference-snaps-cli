package disk

import (
	"fmt"
	"testing"
)

func TestGetDirStatsSnaps(t *testing.T) {
	diskStats, err := GetDirStats("/var/lib/snapd/snaps")
	if err != nil {
		t.Fatalf("GetDirStats() failed: %v", err)
	}

	fmt.Printf("Total: %.0fGB\n", float64(diskStats.Total)/1024/1024/1024)
	fmt.Printf("Used: %.0fGB\n", float64(diskStats.Used)/1024/1024/1024)
	fmt.Printf("Avail: %.0fGB\n", float64(diskStats.Avail)/1024/1024/1024)
}

func TestGetDirStatsRoot(t *testing.T) {
	diskStats, err := GetDirStats("/")
	if err != nil {
		t.Fatalf("GetDirStats() failed: %v", err)
	}

	fmt.Printf("Total: %.0fGB\n", float64(diskStats.Total)/1024/1024/1024)
	fmt.Printf("Used: %.0fGB\n", float64(diskStats.Used)/1024/1024/1024)
	fmt.Printf("Avail: %.0fGB\n", float64(diskStats.Avail)/1024/1024/1024)
}
