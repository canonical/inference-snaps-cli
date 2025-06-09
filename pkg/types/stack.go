package types

type StackSelection struct {
	Stacks   []ScoredStack `json:"stacks"`
	TopStack string        `json:"top-stack"`
}

type ScoredStack struct {
	Stack
	Score      int      `json:"score"`
	Compatible bool     `json:"compatible"`
	Notes      []string `json:"notes,omitempty"`
}

type Stack struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Vendor      string `yaml:"vendor" json:"vendor"`
	Grade       string `yaml:"grade" json:"grade"`

	Devices   StackDevices `yaml:"devices" json:"devices"`
	Memory    *string      `yaml:"memory" json:"memory"`
	DiskSpace *string      `yaml:"disk-space" json:"disk-space"`

	Components     []string  `yaml:"components" json:"components"`
	Configurations StackConf `yaml:"configurations" json:"configurations"`
}

type StackDevices struct {
	Any []StackDevice `yaml:"any" json:"any"`
	All []StackDevice `yaml:"all" json:"all"`
}

type StackDevice struct {
	Type string `yaml:"type" json:"type"`
	Bus  string `yaml:"bus" json:"bus,omitempty"`

	// CPUs
	Architectures []string `yaml:"architectures" json:"architectures,omitempty"`
	VendorId      *string  `yaml:"vendor-id" json:"vendor-id"`
	ModelIds      []string `yaml:"model-id" json:"model-id,omitempty"`
	FamilyIds     []string `yaml:"family-ids" json:"family-ids,omitempty"`
	Flags         []string `yaml:"flags" json:"flags,omitempty"`

	// PCI
	PciDeviceClass *string  `yaml:"pci-device-class" json:"pci-device-class,omitempty"`
	PciVendorId    *string  `yaml:"pci-vendor-id" json:"pci-vendor-id,omitempty"`
	PciDeviceIds   []string `yaml:"pci-device-ids" json:"pci-device-ids,omitempty"`

	// GPU additional properties
	VRam              *string `yaml:"vram" json:"vram,omitempty"`
	ComputeCapability *string `yaml:"compute-capability" json:"compute-capability,omitempty"`
}

type StackConf map[string]interface{}
