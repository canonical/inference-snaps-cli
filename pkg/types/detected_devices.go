package types

type PlatformInfo struct {
	Vendor string `json:"vendor" yaml:"vendor"`
	Name   string `json:"name" yaml:"name"`
	SoC  string `json:"soc,omitempty" yaml:"soc,omitempty"`
}

type DetectedDevice struct {
	Type  string   `json:"type" yaml:"type"`
	Bus   string   `json:"bus" yaml:"bus"`
	PlatformInfo *PlatformInfo `json:"platform_info,omitempty" yaml:"platform_info,omitempty"`
}
