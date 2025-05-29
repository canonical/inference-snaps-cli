package types

type PciPeripheral struct {
	DeviceClass          string  `json:"device_class"`
	ProgrammingInterface *string `json:"programming_interface"`
	VendorId             string  `json:"vendor_id"`
	DeviceId             string  `json:"device_id"`
	SubvendorId          *string `json:"subvendor_id,omitempty"`
	SubdeviceId          *string `json:"subdevice_id,omitempty"`
	VendorName           *string `json:"vendor_name,omitempty"`
	DeviceName           *string `json:"device_name,omitempty"`
	SubvendorName        *string `json:"subvendor_name,omitempty"`
	SubdeviceName        *string `json:"subdevice_name,omitempty"`
}
