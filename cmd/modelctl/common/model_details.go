package common

import (
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
)

type ModelDetails struct {
	Name  string `json:"name" yaml:"name"`
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`

	Description  string   `json:"description,omitempty" yaml:"description,omitempty"`
	ModelCardUrl string   `json:"model-card-url,omitempty" yaml:"model-card-url,omitempty"`
	Format       string   `json:"format,omitempty" yaml:"format,omitempty"`
	Quantization string   `json:"quantization,omitempty" yaml:"quantization,omitempty"`
	Capabilities []string `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`

	DiskSize string `json:"disk-size" yaml:"disk-size"`

	RequiredMemory string `json:"required-memory,omitempty" yaml:"required-memory,omitempty"`
	RequiredVram   string `json:"required-vram,omitempty" yaml:"required-vram,omitempty"`

	Components []string `json:"components,omitempty" yaml:"components,omitempty"`

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
	modelDetails.Components = manifest.Components
	modelDetails.CompatibleEngines = manifest.CompatibilityReport.CompatibleEngines
	modelDetails.Compatible = manifest.CompatibilityReport.Compatible

	// Change disk size to largest possible unit representation
	diskSizeBytes, err := utils.StringToBytes(manifest.DiskSize)
	if err != nil {
		return modelDetails, err
	}
	modelDetails.DiskSize = utils.FmtBytesShort(diskSizeBytes)
	modelDetails.RequiredMemory = utils.FmtBytesShort(manifest.CompatibilityReport.RequiredMemory)
	if manifest.CompatibilityReport.RequiredVram > 0 {
		modelDetails.RequiredVram = utils.FmtBytesShort(manifest.CompatibilityReport.RequiredVram)
	}

	modelDetails.fillIncompatibilityIssues(manifest.CompatibilityReport)
	return modelDetails, nil
}

func (m *ModelDetails) fillIncompatibilityIssues(report models.CompatibilityReport) {
	var issues []string
	if !report.CompatibleMemory {
		if report.TotalVram < report.RequiredVram {
			issues = append(issues, "insufficient VRAM")
		} else if report.TotalMemory < report.RequiredMemory {
			issues = append(issues, "insufficient system memory")
		}
	}

	if !report.CompatibleDisk {
		issues = append(issues, "insufficient disk space")
	}

	m.CompatibilityIssues = issues
}
