package lspci

import (
	"os/exec"
	"strconv"
	"strings"
)

func hostLsPci() ([]byte, error) {
	out, err := exec.Command("lspci", "-vmmnD").Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func parseLsPci(input []byte) ([]PciDevice, error) {
	var devices []PciDevice

	inputString := string(input)
	for _, section := range strings.Split(inputString, "\n\n") {
		var device PciDevice
		for _, line := range strings.Split(section, "\n") {
			key, value, _ := strings.Cut(line, ":\t")

			switch key {
			case "Slot":
				device.Slot = value
			case "Class":
				// e.g. 0x0300 for VGA controller
				if class, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.DeviceClass = uint16(class)
				}
			case "Vendor":
				if vendor, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.VendorId = uint16(vendor)
				}
			case "Device":
				if deviceId, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.DeviceId = uint16(deviceId)
				}
			case "SVendor":
				if subVendorId, err := strconv.ParseUint(value, 16, 16); err == nil {
					subVendorIdUint16 := uint16(subVendorId)
					device.SubVendorId = &subVendorIdUint16
				}
			case "SDevice":
				if subDeviceId, err := strconv.ParseUint(value, 16, 16); err == nil {
					subDeviceIdUint16 := uint16(subDeviceId)
					device.SubDeviceId = &subDeviceIdUint16
				}
			case "ProgIf":
				// e.g. 0x02
				if progIf, err := strconv.ParseUint(value, 16, 8); err == nil {
					progIfUint8 := uint8(progIf)
					device.ProgrammingInterface = &progIfUint8
				}
			}

		}
		devices = append(devices, device)
	}

	return devices, nil
}
