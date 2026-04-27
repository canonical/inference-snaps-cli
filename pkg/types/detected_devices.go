package types

type DeviceMetadata struct {
	VendorName  string `json:"vendor_name" yaml:"vendor_name"`
	ProductName string `json:"product_name,omitempty" yaml:"product_name,omitempty"`
}

type DetectedDevice struct {
	Type     string          `json:"type" yaml:"type"`
	Bus      string          `json:"bus" yaml:"bus"`
	Metadata *DeviceMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}
