package disk

type SystemDirsInfo struct {
	// Root system disk
	Root *DirStats `json:"root"`
	// Snaps read only directory /var/lib/snapd/snaps
	Snaps *DirStats `json:"snaps"`
	// Snaps writable storage /var/snap - not handling this now. We assume the read only directory is on the same file system as the writable directory
	// SnapRW *DirStats
}

type DirStats struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Avail uint64 `json:"avail"`
}
