package types

type PlatformInfo struct {
	Vendor string `json:"vendor" yaml:"vendor"`
	Name   string `json:"name" yaml:"name"`
}

type DetectedDevice struct {
	Type  string   `json:"type" yaml:"type"`
	Bus   string   `json:"bus" yaml:"bus"`
	Nodes []string `json:"nodes,omitempty" yaml:"nodes,omitempty"`
}
