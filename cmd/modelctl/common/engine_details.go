package common

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
)

type EngineDetails struct {
	engines.ScoredManifest `yaml:",inline"`
	Compatible             bool     `yaml:"compatible" json:"compatible"`
	CompatibilityIssues    []string `yaml:"compatibility-issues,omitempty" json:"compatibility-issues,omitempty"`
}

func NewEngineDetails(scoredManifest engines.ScoredManifest) EngineDetails {
	e := EngineDetails{
		ScoredManifest: scoredManifest,
		Compatible:     scoredManifest.CompatibilityReport.EngineCompatible(),
	}
	e.fillIncompatibilityIssues(scoredManifest.CompatibilityReport)
	return e
}

func (e *EngineDetails) fillIncompatibilityIssues(report engines.CompatibilityReport) {
	var issues []string
	if !report.CompatibleMemory {
		var memoryType string
		if report.RequiredMemory.IsVram {
			memoryType = "VRAM"
		} else {
			memoryType = "memory"
		}
		issues = append(issues, fmt.Sprintf("insufficient %s", memoryType))
	}
	if !report.CompatibleDisk {
		issues = append(issues, "insufficient disk space")
	}
	if !report.CompatibleDevices {
		issues = append(issues, "required device not found")
	}
	e.CompatibilityIssues = issues
}
