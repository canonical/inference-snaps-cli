package common

import (
	"fmt"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

func GetIncompatibilityReasons(report engines.CompatibilityReport) []string {
	var reasons []string
	if report.HasMemoryIssue {
		reasons = append(reasons, "insufficient memory")
	}
	if report.HasDiskIssue {
		reasons = append(reasons, "insufficient disk space")
	}
	if report.HasDeviceIssue {
		reasons = append(reasons, "required device not found")
	}
	return reasons
}

func GetVerboseIncompatibilityReasons(report engines.CompatibilityReport) []string {
	var reasons []string
	if report.HasMemoryIssue {
		reasons = append(reasons, fmt.Sprintf("requires %s memory, has %s", utils.FmtBytes(report.RequiredMemory), utils.FmtBytes(report.AvailableMemory)))
	}
	if report.HasDiskIssue {
		reasons = append(reasons, fmt.Sprintf("requires %s disk space, has %s", utils.FmtBytes(report.RequiredDiskSpace), utils.FmtBytes(report.AvailableDiskSpace)))
	}
	if report.HasDeviceIssue {
		if len(report.MissingDevices) > 0 {
			reasons = append(reasons, fmt.Sprintf("required device not found: %s", strings.Join(report.MissingDevices, ", ")))
		} else {
			reasons = append(reasons, "required device not found")
		}
	}
	return reasons
}
