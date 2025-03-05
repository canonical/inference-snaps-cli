package types

type StackSelection struct {
	Stacks   []StackResult `json:"stacks"`
	TopStack string        `json:"top-stack"`
}

type StackResult struct {
	Name       string   `json:"name"`
	Compatible bool     `json:"compatible"`
	Grade      string   `json:"grade"`
	Notes      []string `json:"notes"`
	Size       uint64   `json:"size"` // size on disk
}

type Stack struct {
	Name           string       `yaml:"name"`
	Description    string       `yaml:"description"`
	Maintainer     string       `yaml:"maintainer"`
	Grade          string       `yaml:"grade"`
	Devices        StackDevices `yaml:"devices"`
	Memory         string       `yaml:"memory"`
	DiskSpace      string       `yaml:"disk-space"`
	Components     []string     `yaml:"components"`
	Configurations StackConf    `yaml:"configurations"`
}

type StackDevices struct {
	Any []StackDevice `yaml:"any"`
	All []StackDevice `yaml:"all"`
}

type StackDevice struct {
	Type     string   `yaml:"type"`
	VendorId *string  `yaml:"vendor-id"`
	ModelIDs []string `yaml:"model-ids"`

	// CPUs
	Architectures []string `yaml:"architectures"`
	FamilyIDs     []string `yaml:"family-ids"`
	Flags         []string `yaml:"flags"`

	// GPUs
	Bus               *string `yaml:"bus"`
	ComputeCapability *string `yaml:"compute-capability"`
	MinimumVram       *string `yaml:"vram"` // TODO update key to minimum-vram
}

type StackConf map[string]interface{}
