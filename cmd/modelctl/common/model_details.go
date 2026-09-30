package common

import (
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
)

type ModelDetails struct {
	Name  string `json:"name" yaml:"name"`
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`

	Description  string   `json:"description" yaml:"description"`
	ModelCardUrl string   `json:"model-card-url" yaml:"model-card-url"`
	Format       string   `json:"format" yaml:"format"`
	Quantization string   `json:"quantization" yaml:"quantization"`
	Capabilities []string `json:"capabilities" yaml:"capabilities"`

	DiskSize string `json:"disk-size" yaml:"disk-size"`

	KVCacheSize string `json:"kv-cache-size,omitempty" yaml:"kv-cache-size,omitempty"`

	Components []string `json:"components" yaml:"components"`

	CompatibleEngines []string `json:"compatible-engines,omitempty" yaml:"compatible-engines,omitempty"`

	Compatible          bool     `json:"compatible" yaml:"compatible"`
	CompatibilityIssues []string `json:"compatibility-issues,omitempty" yaml:"compatibility-issues,omitempty"`
}

func NewModelDetails(manifest *models.ScoredManifest) (ModelDetails, error) {
	var modelDetails ModelDetails
	modelDetails.Name = manifest.Name
	modelDetails.Alias = manifest.Alias
	modelDetails.Description = manifest.Description
	modelDetails.ModelCardUrl = manifest.ModelCardUrl
	modelDetails.Format = manifest.Format
	modelDetails.Quantization = manifest.Quantization
	modelDetails.Capabilities = manifest.Capabilities
	modelDetails.KVCacheSize = manifest.KVCacheSize
	modelDetails.Components = manifest.Components
	modelDetails.CompatibleEngines = manifest.CompatibilityReport.CompatibleEngines
	modelDetails.Compatible = manifest.CompatibilityReport.Compatible

	// Change disk size to largest possible unit representation
	diskSizeBytes, err := utils.StringToBytes(manifest.DiskSize)
	if err != nil {
		return modelDetails, ErrInsufficientDiskSpaceForModel
	}
	modelDetails.DiskSize = utils.FmtBytesShort(diskSizeBytes)

	modelDetails.fillIncompatibilityIssues(manifest.CompatibilityReport)
	return modelDetails, nil
}

func (m *ModelDetails) fillIncompatibilityIssues(report models.CompatibilityReport) {
	var issues []string
	if !report.CompatibleMemory {
		issues = append(issues, "insufficient memory")
	}
	if !report.CompatibleDisk {
		issues = append(issues, "insufficient disk space")
	}

	m.CompatibilityIssues = issues
}
