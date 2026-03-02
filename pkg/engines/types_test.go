package engines

import (
	"testing"
)

func TestGetCompatibilityIssues(t *testing.T) {
	var CompatibilityIssues []CompatibilityIssue
	CompatibilityIssues = append(CompatibilityIssues, DeviceCompatibilityIssue{})
	CompatibilityIssues = append(CompatibilityIssues, MemoryCompatibilityIssue{
		RequiredMemory:  8 * 1024 * 1024 * 1024, // 8 GiB
		AvailableMemory: 4 * 1024 * 1024 * 1024, // 4 GiB
	})
	CompatibilityIssues = append(CompatibilityIssues, DiskCompatibilityIssue{
		RequiredSpace:  100 * 1024 * 1024 * 1024, // 100 GiB
		AvailableSpace: 50 * 1024 * 1024 * 1024,  // 50 GiB
	})

	if len(CompatibilityIssues) != 3 {
		t.Fatalf("Expected to have 3 compatibility issues in the slice, got: %d", len(CompatibilityIssues))
	}

	var reasons []string
	for _, issue := range CompatibilityIssues {
		reasons = append(reasons, issue.GetReason())
	}

	expectedReasons := []string{
		"required device not found",
		"insufficient memory",
		"insufficient disk space",
	}

	for i, reason := range reasons {
		if reason != expectedReasons[i] {
			t.Errorf("Expected reason: %s, got: %s", expectedReasons[i], reason)
		}
	}
}

func TestGetVerboseCompatibilityIssues(t *testing.T) {
	var CompatibilityIssues []CompatibilityIssue
	CompatibilityIssues = append(CompatibilityIssues, DeviceCompatibilityIssue{})
	CompatibilityIssues = append(CompatibilityIssues, MemoryCompatibilityIssue{
		RequiredMemory:  8 * 1024 * 1024 * 1024, // 8 GiB
		AvailableMemory: 4 * 1024 * 1024 * 1024, // 4 GiB
	})
	CompatibilityIssues = append(CompatibilityIssues, DiskCompatibilityIssue{
		RequiredSpace:  100 * 1024 * 1024 * 1024, // 100 GiB
		AvailableSpace: 50 * 1024 * 1024 * 1024,  // 50 GiB
	})

	if len(CompatibilityIssues) != 3 {
		t.Fatalf("Expected to have 3 compatibility issues in the slice, got: %d", len(CompatibilityIssues))
	}

	var reasons []string
	for _, issue := range CompatibilityIssues {
		reasons = append(reasons, issue.GetVerboseReason())
	}

	expectedReasons := []string{
		"required device not found",
		"requires 8.0GiB memory, has 4.0GiB",
		"requires 100.0GiB disk space, has 50.0GiB",
	}

	for i, reason := range reasons {
		if reason != expectedReasons[i] {
			t.Errorf("Expected reason: %s, got: %s", expectedReasons[i], reason)
		}
	}
}
