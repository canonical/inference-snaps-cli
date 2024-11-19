package gpu

type Display struct {
	Vendor  string `json:"vendor"`
	Product string `json:"product"`
}

type Gpu struct {
	VendorId    string  `json:"vendor_id"`
	DeviceId    string  `json:"device_id"`
	SubVendorId *string `json:"sub_vendor_id,omitempty"`
	SubDeviceId *string `json:"sub_device_id,omitempty"`
	VRam        *uint64 `json:"vram,omitempty"`
}
