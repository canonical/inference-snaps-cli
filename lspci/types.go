package lspci

type PciDevice struct {
	Slot                 string  `json:"slot"`
	DeviceClass          uint16  `json:"device_class"`
	ProgrammingInterface *uint8  `json:"programming_interface"`
	VendorId             uint16  `json:"vendor_id"`
	VendorName           *string `json:"vendor_name,omitempty"`
	DeviceId             uint16  `json:"device_id"`
	DeviceName           *string `json:"device_name,omitempty"`
	SubVendorId          *uint16 `json:"sub_vendor_id,omitempty"`
	SubDeviceId          *uint16 `json:"sub_device_id,omitempty"`
}
