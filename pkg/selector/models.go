package selector

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/v2/pkg/constants"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/canonical/lscompute/pkg/machine"
)

func SelectModel(modelOptions []string, preferredModel string, manifests map[string]models.Manifest, machineInfo *machine.MachineInfo) (string, error) {
	availableDiskSpace, err := availableDiskSpace(machineInfo)
	if err != nil {
		return "", err
	}

	if preferredModel != "" {
		fitsPreferred, _, err := modelFitsDiskSpace(preferredModel, availableDiskSpace, manifests)
		if err != nil {
			return "", err
		}
		if fitsPreferred {
			return preferredModel, nil
		}
	}

	selected := ""
	var selectedSize uint64
	smallestModel := ""
	smallestSize := ^uint64(0)
	for _, modelID := range modelOptions {
		fitsModel, size, err := modelFitsDiskSpace(modelID, availableDiskSpace, manifests)
		if err != nil {
			return "", err
		}
		if size < smallestSize {
			smallestModel = modelID
			smallestSize = size
		}
		if fitsModel && (selected == "" || size > selectedSize) {
			selected = modelID
			selectedSize = size
		}
	}

	if selected != "" {
		return selected, nil
	}
	if smallestModel != "" {
		return smallestModel, nil
	}
	return "", fmt.Errorf("no model options available")
}

func availableDiskSpace(machineInfo *machine.MachineInfo) (uint64, error) {

	disk, ok := machineInfo.Disk[constants.SnapStoragePath]
	if !ok {
		return 0, fmt.Errorf("disk information unavailable for %s", constants.SnapStoragePath)
	}
	return disk.Avail, nil
}

func modelFitsDiskSpace(modelID string, availableDiskSpace uint64, manifests map[string]models.Manifest) (bool, uint64, error) {
	manifest, ok := manifests[modelID]
	if !ok {
		return false, 0, fmt.Errorf("model manifest not found: %s", modelID)
	}
	size, err := utils.StringToBytes(manifest.DiskSize)
	if err != nil {
		return false, 0, fmt.Errorf("parsing disk size for model %q: %w", modelID, err)
	}
	return size <= availableDiskSpace, size, nil
}
