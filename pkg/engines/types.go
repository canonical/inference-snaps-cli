package engines

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"gopkg.in/yaml.v3"
)

type SpaceCompatibilityStatus struct {
	Compatible     bool
	RequiredSpace  uint64
	AvailableSpace uint64
}

type DeviceCompatibilityStatus struct {
	Compatible bool
}

type CompatibilityIssues struct {
	Memory SpaceCompatibilityStatus
	Disk   SpaceCompatibilityStatus
	Device DeviceCompatibilityStatus
}

type ScoredManifest struct {
	Manifest            `yaml:",inline"`
	Score               int                 `yaml:"score" json:"score"`
	Compatible          bool                `yaml:"compatible" json:"compatible"`
	CompatibilityIssues CompatibilityIssues `yaml:"compatibility-issues,omitempty" json:"compatibility-issues,omitempty"`
}

type Manifest struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Vendor      string `yaml:"vendor" json:"vendor"`
	Grade       string `yaml:"grade" json:"grade"`

	Devices   Devices `yaml:"devices" json:"devices"`
	Memory    *string `yaml:"memory,omitempty" json:"memory"`
	DiskSpace *string `yaml:"disk-space,omitempty" json:"disk-space"`

	Components     []string       `yaml:"components" json:"components"`
	Configurations Configurations `yaml:"configurations" json:"configurations"`
}

type Devices struct {
	Anyof []Device `yaml:"anyof,omitempty" json:"anyof"`
	Allof []Device `yaml:"allof,omitempty" json:"allof"`
}

type Device struct {
	// General
	Type string `yaml:"type,omitempty" json:"type,omitempty"` // cpu, gpu, npu or nil
	Bus  string `yaml:"bus,omitempty" json:"bus,omitempty"`   // pci, usb or nil

	// CPUs
	Architecture *string `yaml:"architecture,omitempty" json:"architecture,omitempty"`

	// CPU x86
	ManufacturerId *string  `yaml:"manufacturer-id,omitempty" json:"manufacturer-id,omitempty"`
	Flags          []string `yaml:"flags,omitempty" json:"flags,omitempty"`

	// CPU arm64
	ImplementerId *types.HexInt `yaml:"implementer-id,omitempty" json:"implementer-id,omitempty"`
	PartNumber    *types.HexInt `yaml:"part-number,omitempty" json:"part-number,omitempty"`
	Features      []string      `yaml:"features,omitempty" json:"features,omitempty"`

	// PCI
	VendorId *types.HexInt `yaml:"vendor-id,omitempty" json:"vendor-id,omitempty"`
	DeviceId *types.HexInt `yaml:"device-id,omitempty" json:"device-id,omitempty"`

	// GPU additional properties
	VRam              *string `yaml:"vram,omitempty" json:"vram,omitempty"`
	ComputeCapability *string `yaml:"compute-capability,omitempty" json:"compute-capability,omitempty"`

	// NPU
	// no additional properties for now

	// Drivers
	SnapConnections []string `yaml:"snap-connections,omitempty" json:"snap-connections,omitempty"`

	// General
	CompatibilityIssues []string `yaml:"compatibility-issues,omitempty" json:"compatibility-issues,omitempty"`
}

type Configurations map[string]interface{}

func (c CompatibilityIssues) GetReasons() []string {
	// Returns an array of strings with all the compatibility issues

	var reasons []string

	if !c.Memory.Compatible {
		reasons = append(reasons, "insufficient memory")
	}

	if !c.Disk.Compatible {
		reasons = append(reasons, "insufficient disk space")
	}

	if !c.Device.Compatible {
		reasons = append(reasons, "required device not found")
	}

	return reasons
}

func (c CompatibilityIssues) GetVerboseReasons() []string {
	// Returns an array of strings with detailed compatibility issues

	var reasons []string

	if !c.Memory.Compatible {
		reason := fmt.Sprintf("requires %s memory, has %s", utils.FmtBytes(c.Memory.RequiredSpace), utils.FmtBytes(c.Memory.AvailableSpace))
		reasons = append(reasons, reason)
	}

	if !c.Disk.Compatible {
		reason := fmt.Sprintf("requires %s disk space, has %s", utils.FmtBytes(c.Disk.RequiredSpace), utils.FmtBytes(c.Disk.AvailableSpace))
		reasons = append(reasons, reason)
	}

	if !c.Device.Compatible {
		reasons = append(reasons, "required device not found")
	}

	return reasons
}

func (c CompatibilityIssues) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(c.GetReasons())
}

func (c CompatibilityIssues) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.GetReasons())
}
