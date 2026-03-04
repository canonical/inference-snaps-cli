package common

import (
	"fmt"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

type EngineDetails struct {
	engines.ScoredManifest
	Compatible          bool     `yaml:"compatible" json:"compatible"`
	CompatibilityIssues []string `yaml:"compatibility-issues,omitempty" json:"compatibility-issues,omitempty"`
}

func NewEngineDetails(scoredManifest engines.ScoredManifest) EngineDetails {
	return EngineDetails{
		ScoredManifest:      scoredManifest,
		Compatible:          scoredManifest.CompatibilityReport.IsCompatible(),
		CompatibilityIssues: getIncompatibilityReasons(scoredManifest.CompatibilityReport),
	}
}

func NewVerboseEngineDetails(scoredManifest engines.ScoredManifest) EngineDetails {
	return EngineDetails{
		ScoredManifest:      scoredManifest,
		Compatible:          scoredManifest.CompatibilityReport.IsCompatible(),
		CompatibilityIssues: getVerboseIncompatibilityReasons(scoredManifest.CompatibilityReport),
	}
}

func getIncompatibilityReasons(report engines.CompatibilityReport) []string {
	var reasons []string
	if !report.IsMemoryCompatible {
		reasons = append(reasons, "insufficient memory")
	}
	if !report.IsDiskCompatible {
		reasons = append(reasons, "insufficient disk space")
	}
	if !report.IsDeviceCompatible {
		reasons = append(reasons, "required device not found")
	}
	return reasons
}

func getVerboseIncompatibilityReasons(report engines.CompatibilityReport) []string {
	var reasons []string
	if !report.IsMemoryCompatible {
		reasons = append(reasons, fmt.Sprintf("requires %s memory, has %s", utils.FmtBytes(report.RequiredMemory), utils.FmtBytes(report.AvailableMemory)))
	}
	if !report.IsDiskCompatible {
		reasons = append(reasons, fmt.Sprintf("requires %s disk space, has %s", utils.FmtBytes(report.RequiredDiskSpace), utils.FmtBytes(report.AvailableDiskSpace)))
	}
	if !report.IsDeviceCompatible {
		if len(report.MissingDevices) > 0 {
			reasons = append(reasons, fmt.Sprintf("required device not found: %s", strings.Join(report.MissingDevices, ", ")))
		} else {
			reasons = append(reasons, "required device not found")
		}
	}
	return reasons
}
