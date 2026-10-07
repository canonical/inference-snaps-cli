package common

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/canonical/inference-snaps-cli/v2/pkg/constants"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	selectorpci "github.com/canonical/inference-snaps-cli/v2/pkg/selector/pci"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
)

func GetModelManifestByNameOrAlias(ctx *Context, modelName string, machine *machine.Machine) (*models.ScoredManifest, error) {
	if modelName == "" {
		return nil, fmt.Errorf("model name must not be empty")
	}

	activeEngine, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}
	if activeEngine == "" {
		return nil, ErrNoActiveEngine
	}

	engineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngine)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingEngineManifest, err)
	}

	allModelManifests, err := models.LoadManifests(ctx.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingModelManifests, err)
	}

	// Consider the active engine's models first
	var manifest *models.Manifest
	for _, modelManifest := range allModelManifests {
		if slices.Contains(engineManifest.Model.Options, modelManifest.Name) &&
			(modelManifest.Name == modelName || modelManifest.Alias == modelName) {
			manifest = &modelManifest
			break
		}
	}

	// If the provided name is not one of the active engine's models, warn the user it is not compatible
	if manifest == nil {
		for _, modelManifest := range allModelManifests {
			if modelManifest.Name == modelName || modelManifest.Alias == modelName {
				return nil, fmt.Errorf("model %q is not compatible with the active engine", modelName)
			}
		}
	}

	if manifest == nil {
		return nil, fmt.Errorf("model %q does not exist", modelName)
	}

	scoredManifest, err := GetScoredModel(ctx, *engineManifest, *manifest, machine)
	if err != nil {
		return nil, err
	}

	return &scoredManifest, nil
}

func GetCompatibleEnginesByModelName(ctx *Context, modelName string) ([]string, error) {
	allEngineManifests, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingEngineManifest, err)
	}

	compatibleEngines := []string{}
	for _, engineManifest := range allEngineManifests {
		if slices.Contains(engineManifest.Model.Options, modelName) {
			compatibleEngines = append(compatibleEngines, engineManifest.Name)
		}
	}
	return compatibleEngines, nil
}

func ModelStatus(ctx *Context) (map[string]string, error) {
	activeModelId, err := ctx.Cache.GetActiveModel()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveModel, err)
	}

	if activeModelId == "" {
		return nil, ErrNoActiveModel
	}

	activeModelManifest, err := models.LoadManifest(ctx.ModelsDir, activeModelId)
	if err != nil {
		return nil, fmt.Errorf("loading model manifest: %v", err)
	}

	status := make(map[string]string)
	status["name"] = activeModelManifest.Name

	return status, nil
}

func GetAllModels(ctx *Context, machine *machine.Machine) ([]ModelDetails, error) {
	allModelManifests, err := models.LoadManifests(ctx.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingModelManifests, err)
	}

	manifestsByName := make(map[string]models.Manifest, len(allModelManifests))
	for _, manifest := range allModelManifests {
		manifestsByName[manifest.Name] = manifest
	}

	allEngineManifests, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingEngineManifest, err)
	}

	activeEngine, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}

	var allScoredModelManifests []models.ScoredManifest
	for _, engineManifest := range allEngineManifests {
		if engineManifest.Name == activeEngine {
			allScoredModelManifests, err = ScoreModelsAgainstEngine(ctx, engineManifest, manifestsByName, machine)
			if err != nil {
				return nil, fmt.Errorf("scoring models: %v", err)
			}
			break
		}
	}
	var allModelsDetails []ModelDetails
	for _, manifest := range allScoredModelManifests {
		modelDetails, err := NewModelDetails(&manifest)
		if err != nil {
			return nil, fmt.Errorf("creating model details for %s: %w", manifest.Name, err)
		}
		allModelsDetails = append(allModelsDetails, modelDetails)
	}
	slices.SortFunc(allModelsDetails, func(a, b ModelDetails) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return allModelsDetails, nil
}

func GetScoredModel(ctx *Context, engineManifest engines.Manifest, modelManifest models.Manifest, machine *machine.Machine) (models.ScoredManifest, error) {
	manifests := map[string]models.Manifest{
		modelManifest.Name: modelManifest,
	}
	scoredModels, err := ScoreModelsAgainstEngine(ctx, engineManifest, manifests, machine)
	if err != nil {
		return models.ScoredManifest{}, err
	}
	for _, scoredModel := range scoredModels {
		if scoredModel.Manifest.Name == modelManifest.Name {
			return scoredModel, nil
		}
	}
	return models.ScoredManifest{}, fmt.Errorf("model not found: %s", modelManifest.Name)
}

func ScoreModelsAgainstEngine(ctx *Context, engineManifest engines.Manifest, modelManifests map[string]models.Manifest, machine *machine.Machine) ([]models.ScoredManifest, error) {
	availableDiskSpace, err := availableDiskSpace(machine)
	if err != nil {
		return nil, err
	}

	availableMemory, err := availableMemory(machine, engineManifest)
	if err != nil {
		return nil, err
	}

	var runtimeMemory uint64 = 200 * 1024 * 1024 // 200 MB
	if engineManifest.Runtime != "" {
		runtimeManifest, err := runtimes.LoadManifest(ctx.RuntimesDir, engineManifest.Runtime)
		if err != nil {
			return nil, fmt.Errorf("loading runtime manifest: %v", err)
		}
		if runtimeManifest.RequiredMemory != "" {
			runtimeMemory, err = utils.StringToBytes(runtimeManifest.RequiredMemory)
			if err != nil {
				return nil, fmt.Errorf("parsing runtime memory: %v", err)
			}
		}
	}

	engineManifests, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("loading engine manifests: %w", err)
	}

	scoredModels := make([]models.ScoredManifest, 0, len(modelManifests))
	for modelID := range modelManifests {
		manifest := modelManifests[modelID]

		size, err := utils.StringToBytes(manifest.DiskSize)
		if err != nil {
			return nil, fmt.Errorf("parsing disk size for model %q: %w", modelID, err)
		}

		var requiredMemory uint64
		if manifest.RequiredMemory != "" {
			requiredMemory, err = utils.StringToBytes(manifest.RequiredMemory)
			if err != nil {
				return nil, fmt.Errorf("parsing required memory for model %q: %w", modelID, err)
			}
		} else {
			requiredMemory = size
		}

		// 2GB for squashfs and OS
		requiredMemory, ok := utils.CheckedAddUint64(requiredMemory, runtimeMemory, 2*1024*1024*1024)
		if !ok {
			return nil, fmt.Errorf("required memory for model %q overflows uint64", modelID)
		}

		report := models.CompatibilityReport{
			CompatibleDisk:     size <= availableDiskSpace,
			RequiredDiskSpace:  size,
			AvailableDiskSpace: availableDiskSpace,
			RequiredMemory:     requiredMemory,
			AvailableMemory:    availableMemory,
		}
		report.CompatibleMemory = report.RequiredMemory <= availableMemory
		report.Compatible = report.CompatibleDisk && report.CompatibleMemory

		score := uint64(0)
		if report.CompatibleDisk && report.CompatibleMemory {
			score, ok = utils.CheckedAddUint64(report.RequiredDiskSpace, report.RequiredMemory)
			if !ok {
				score = math.MaxUint64
			}
		}
		scoredModelManifest := models.ScoredManifest{
			Manifest:            manifest,
			Score:               score,
			CompatibilityReport: report,
		}
		compatibleEngines := []string{}
		for _, engineManifest := range engineManifests {
			if slices.Contains(engineManifest.Model.Options, manifest.Name) {
				compatibleEngines = append(compatibleEngines, engineManifest.Name)
			}
		}
		scoredModelManifest.CompatibilityReport.CompatibleEngines = compatibleEngines
		scoredModels = append(scoredModels, scoredModelManifest)
	}
	return scoredModels, nil
}

// SelectModel picks a model that fits the available disk: the preferred one if it fits, otherwise the largest fitting model.
// Returns ErrInsufficientDiskSpaceForModel when none fit.
func SelectModel(ctx *Context, engineManifest engines.Manifest, preferredModel string, machine *machine.Machine) (string, []models.ScoredManifest, error) {
	modelManifests, err := models.LoadManifests(ctx.ModelsDir)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", LoadingModelManifests, err)
	}
	modelManifestsByName := make(map[string]models.Manifest)
	for _, manifest := range modelManifests {
		modelManifestsByName[manifest.Name] = manifest
	}

	filteredManifestsByName := make(map[string]models.Manifest)
	for _, model := range engineManifest.Model.Options {
		manifest, ok := modelManifestsByName[model]
		if !ok {
			return "", nil, fmt.Errorf("model manifest not found: %s", model)
		}
		filteredManifestsByName[model] = manifest
	}
	modelManifestsByName = filteredManifestsByName

	scoredModels, err := ScoreModelsAgainstEngine(ctx, engineManifest, modelManifestsByName, machine)
	if err != nil {
		return "", nil, err
	}
	// Keep the engine manifest's order; scoring iterates a map.
	slices.SortFunc(scoredModels, func(a, b models.ScoredManifest) int {
		return cmp.Compare(slices.Index(engineManifest.Model.Options, a.Name), slices.Index(engineManifest.Model.Options, b.Name))
	})

	if preferredModel != "" {
		for _, model := range scoredModels {
			if model.Name == preferredModel && model.CompatibilityReport.Compatible {
				return preferredModel, scoredModels, nil
			}
		}
	}

	selected := ""
	selectedScore := uint64(0)
	for _, model := range scoredModels {
		if model.CompatibilityReport.Compatible && (selected == "" || model.Score > selectedScore) {
			selected = model.Name
			selectedScore = model.Score
		}
	}

	if selected != "" {
		return selected, scoredModels, nil
	}

	return "", scoredModels, ErrInsufficientDiskSpaceForModel
}

func availableDiskSpace(machine *machine.Machine) (uint64, error) {
	for _, disk := range machine.Disk {
		if disk.Path == constants.SnapStoragePath {
			return disk.Available, nil
		}
	}
	return 0, fmt.Errorf("disk information unavailable for %s", constants.SnapStoragePath)
}

func availableMemory(machine *machine.Machine, engineManifest engines.Manifest) (uint64, error) {
	systemMemory := machine.Memory.TotalRam

	var gpuDevices []engines.Device
	for _, device := range slices.Concat(engineManifest.Devices.Allof, engineManifest.Devices.Anyof) {
		if device.Type == "gpu" && (device.Bus == "" || device.Bus == "pci") {
			gpuDevices = append(gpuDevices, device)
		}
	}
	if len(gpuDevices) == 0 {
		return systemMemory, nil
	}

	// Use the host GPU that best satisfies the engine, mirroring engine selection scoring
	var bestDevice *pci.Device
	bestScore := 0
	for _, gpuDevice := range gpuDevices {
		hostDevice, score := selectorpci.BestMatch(gpuDevice, machine)
		if hostDevice != nil && score > bestScore {
			bestDevice = hostDevice
			bestScore = score
		}
	}
	if bestDevice == nil {
		return systemMemory, nil
	}

	vram, found := bestDevice.AdditionalProperties["vram"]
	if !found || vram == "[N/A]" {
		// assuming unified memory
		return systemMemory, nil
	}
	vramVal, err := strconv.ParseUint(vram, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing vram %q of pci device %s: %w", vram, bestDevice.Slot, err)
	}
	return vramVal, nil
}
