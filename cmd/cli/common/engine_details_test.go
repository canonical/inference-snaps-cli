package common

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
)

func TestGetIncompatibilityReasons(t *testing.T) {
	var compatibilityReport engines.CompatibilityReport = engines.CompatibilityReport{
		HasMemoryIssue:     true,
		RequiredMemory:     8 * 1024 * 1024 * 1024, // 8 GiB
		AvailableMemory:    4 * 1024 * 1024 * 1024, // 4 GiB
		HasDiskIssue:       true,
		RequiredDiskSpace:  100 * 1024 * 1024 * 1024, // 100 GiB
		AvailableDiskSpace: 50 * 1024 * 1024 * 1024,  // 50 GiB
		HasDeviceIssue:     true,
		MissingDevices:     []string{"cpu"},
	}

	expectedReasons := []string{
		"insufficient memory",
		"insufficient disk space",
		"required device not found",
	}

	var actualReasons []string = getIncompatibilityReasons(compatibilityReport)

	if len(actualReasons) != len(expectedReasons) {
		t.Fatalf("Expected to have %d compatibility issues, got: %d", len(expectedReasons), len(actualReasons))
	}

	for i, reason := range actualReasons {
		if reason != expectedReasons[i] {
			t.Errorf("Expected reason: %s, got: %s", expectedReasons[i], reason)
		}
	}
}

func TestGetVerboseIncompatibilityReasons(t *testing.T) {
	var compatibilityReport engines.CompatibilityReport = engines.CompatibilityReport{
		HasMemoryIssue:     true,
		RequiredMemory:     8 * 1024 * 1024 * 1024, // 8 GiB
		AvailableMemory:    4 * 1024 * 1024 * 1024, // 4 GiB
		HasDiskIssue:       true,
		RequiredDiskSpace:  100 * 1024 * 1024 * 1024, // 100 GiB
		AvailableDiskSpace: 50 * 1024 * 1024 * 1024,  // 50 GiB
		HasDeviceIssue:     true,
		MissingDevices:     []string{"cpu"},
	}

	expectedReasons := []string{
		"requires 8.0GiB memory, has 4.0GiB",
		"requires 100.0GiB disk space, has 50.0GiB",
		"required device not found: cpu",
	}

	var actualReasons []string = getVerboseIncompatibilityReasons(compatibilityReport)

	if len(actualReasons) != len(expectedReasons) {
		t.Fatalf("Expected to have %d compatibility issues, got: %d", len(expectedReasons), len(actualReasons))
	}

	for i, reason := range actualReasons {
		if reason != expectedReasons[i] {
			t.Errorf("Expected reason: %s, got: %s", expectedReasons[i], reason)
		}
	}
}
