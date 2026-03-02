package engines

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

type CompatibilityIssue interface {
	GetReason() string
	GetVerboseReason() string
	MarshalYAML() (any, error)
	MarshalJSON() ([]byte, error)
}

type MemoryCompatibilityIssue struct {
	RequiredMemory  uint64
	AvailableMemory uint64
}

type DiskCompatibilityIssue struct {
	RequiredSpace  uint64
	AvailableSpace uint64
}

type DeviceCompatibilityIssue struct{}

func (issue MemoryCompatibilityIssue) GetReason() string {
	return "insufficient memory"
}

func (issue MemoryCompatibilityIssue) GetVerboseReason() string {
	return fmt.Sprintf("requires %s memory, has %s", utils.FmtBytes(issue.RequiredMemory), utils.FmtBytes(issue.AvailableMemory))
}

func (issue DiskCompatibilityIssue) GetReason() string {
	return "insufficient disk space"
}

func (issue DiskCompatibilityIssue) GetVerboseReason() string {
	return fmt.Sprintf("requires %s disk space, has %s", utils.FmtBytes(issue.RequiredSpace), utils.FmtBytes(issue.AvailableSpace))
}

func (issue DeviceCompatibilityIssue) GetReason() string {
	return "required device not found"
}

func (issue DeviceCompatibilityIssue) GetVerboseReason() string {
	return "required device not found"
}

func (issue MemoryCompatibilityIssue) MarshalYAML() (any, error) {
	return issue.GetReason(), nil
}

func (issue MemoryCompatibilityIssue) MarshalJSON() ([]byte, error) {
	return json.Marshal(issue.GetReason())
}

func (issue DiskCompatibilityIssue) MarshalYAML() (any, error) {
	return issue.GetReason(), nil
}

func (issue DiskCompatibilityIssue) MarshalJSON() ([]byte, error) {
	return json.Marshal(issue.GetReason())
}

func (issue DeviceCompatibilityIssue) MarshalYAML() (any, error) {
	return issue.GetReason(), nil
}

func (issue DeviceCompatibilityIssue) MarshalJSON() ([]byte, error) {
	return json.Marshal(issue.GetReason())
}
