package engines

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func (device Device) validateBus(extraFields []string) error {
	switch device.Bus {
	case "pci":
		return device.validatePci(extraFields)
	case "usb":
		return device.validateUsb(extraFields)
	case "fastrpc":
		return device.validateFastRpc(extraFields)
	case "": // default to pci bus
		return device.validatePci(extraFields)
	default:
		return fmt.Errorf("invalid bus: %v", device.Bus)
	}
}

func (device Device) validateUsb(extraFields []string) error {
	return fmt.Errorf("usb: device validation not implemented")
}

func (device Device) validatePci(extraFields []string) error {
	validFields := []string{
		"Type",
		"Bus",
		"VendorId",
		"DeviceId",
		"SnapConnections",
	}
	validFields = append(validFields, extraFields...)

	t := reflect.TypeOf(device)
	v := reflect.ValueOf(device)

	// Check fields with values against allow list
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.FieldByName(fieldName)
		if fieldValue.IsValid() && !fieldValue.IsZero() {
			if !slices.Contains(validFields, fieldName) {
				return fmt.Errorf("pci: invalid field: %s", fieldName)
			}
		}
	}

	return nil
}

func (device Device) validateFastRpc(extraFields []string) error {
	if device.Type != "" && device.Type != "npu" {
		return fmt.Errorf("fastrpc bus only supports npu devices")
	}
	if device.Domain != nil {
		switch strings.ToLower(strings.TrimSpace(*device.Domain)) {
		case "adsp", "mdsp", "sdsp", "cdsp", "gdsp":
		default:
			return fmt.Errorf("fastrpc: invalid domain: %s", *device.Domain)
		}
		if device.Type == "npu" && !strings.EqualFold(strings.TrimSpace(*device.Domain), "cdsp") {
			return fmt.Errorf("fastrpc: npu devices require the cdsp domain")
		}
	}

	validFields := []string{
		"Type",
		"Bus",
		"Domain",
		"SnapConnections",
	}
	validFields = append(validFields, extraFields...)

	t := reflect.TypeOf(device)
	v := reflect.ValueOf(device)

	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.FieldByName(fieldName)
		if fieldValue.IsValid() && !fieldValue.IsZero() {
			if !slices.Contains(validFields, fieldName) {
				return fmt.Errorf("fastrpc: invalid field: %s", fieldName)
			}
		}
	}

	return nil
}
