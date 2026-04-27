package types

type DeviceMetadata struct {
	Vendor     string `json:"vendor" yaml:"vendor"`
	Name       string `json:"name" yaml:"name"`
	VendorName string `json:"vendor_name,omitempty" yaml:"vendor_name,omitempty"`
}

type DetectedDevice struct {
	Type     string          `json:"type" yaml:"type"`
	Bus      string          `json:"bus" yaml:"bus"`
	Metadata *DeviceMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}
