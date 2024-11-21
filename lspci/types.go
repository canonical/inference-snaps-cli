package lspci

type PciDevice struct {
	Slot                 string  `json:"slot"`
	DeviceClass          uint16  `json:"device_class"`
	ProgrammingInterface *uint8  `json:"programming_interface"`
	VendorId             uint16  `json:"vendor_id"`
	VendorName           *string `json:"vendor_name,omitempty"`
	DeviceId             uint16  `json:"device_id"`
	DeviceName           *string `json:"device_name,omitempty"`
	SubvendorId          *uint16 `json:"subvendor_id,omitempty"`
	SubvendorName        *string `json:"subvendor_name,omitempty"`
	SubdeviceId          *uint16 `json:"subdevice_id,omitempty"`
	SubdeviceName        *string `json:"subdevice_name,omitempty"`
}
