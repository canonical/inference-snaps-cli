package gpu

type Display struct {
	Vendor  string `json:"vendor"`
	Product string `json:"product"`
}

type Gpu struct {
	VendorId      string  `json:"vendor_id"`
	VendorName    *string `json:"vendor_name,omitempty"`
	DeviceId      string  `json:"device_id"`
	DeviceName    *string `json:"device_name,omitempty"`
	SubVendorId   *string `json:"sub_vendor_id,omitempty"`
	SubVendorName *string `json:"sub_vendor_name,omitempty"`
	SubDeviceId   *string `json:"sub_device_id,omitempty"`
	SubDeviceName *string `json:"sub_device_name,omitempty"`
	VRam          *uint64 `json:"vram,omitempty"`
}
