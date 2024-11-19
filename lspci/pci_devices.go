package lspci

import (
	"fmt"

	"github.com/jaypipes/pcidb"
)

var (
	pciDb *pcidb.PCIDB
)

func PciDevices(friendlyNames bool) ([]PciDevice, error) {

	hostLsPci, err := hostLsPci()
	if err != nil {
		return nil, err
	}
	devices, err := ParseLsPci(hostLsPci, friendlyNames)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func lookupFriendlyName(device *PciDevice) {
	if pciDb == nil {
		return
	}

	vendorIdString := fmt.Sprintf("%04x", device.VendorId)
	deviceIdString := fmt.Sprintf("%04x", device.DeviceId)

	subVendorIdString := ""
	if device.SubDeviceId != nil {
		subVendorIdString = fmt.Sprintf("%04x", *device.SubVendorId)
	}

	subDeviceIdString := ""
	if device.SubDeviceId != nil {
		subDeviceIdString = fmt.Sprintf("%04x", *device.SubDeviceId)
	}

	for _, vendor := range pciDb.Vendors {
		if vendor.ID == vendorIdString {
			vendorName := vendor.Name
			device.VendorName = &vendorName

			for _, product := range vendor.Products {
				if product.ID == deviceIdString {
					productName := product.Name
					device.DeviceName = &productName

					// Look up subDevice name from subsystem list
					if device.SubDeviceId != nil {
						for _, subSystem := range product.Subsystems {
							if subSystem.ID == subDeviceIdString {
								subSystemName := subSystem.Name
								device.SubDeviceName = &subSystemName
							}
						}
					}
				}
			}
		}

		// Look up SubVendor name from main vendor list
		if device.SubVendorId != nil && vendor.ID == subVendorIdString {
			vendorName := vendor.Name
			device.SubVendorName = &vendorName
		}
	}
}
