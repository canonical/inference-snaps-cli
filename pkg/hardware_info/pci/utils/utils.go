package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// SlotAddress represents the components of a PCI slot address.
type SlotAddress struct {
	Domain   uint32
	Bus      uint32
	Device   uint32
	Function uint32
}

// ParsePCISlot takes a slot string (e.g., "0000:01:00.0" or "01:00.0") and returns a SlotAddress struct.
func ParsePCISlot(slot string) (SlotAddress, error) {
	var device SlotAddress

	// Split the string into parts by the colon delimiter
	parts := strings.Split(slot, ":")

	if len(parts) < 2 || len(parts) > 3 {
		return SlotAddress{}, fmt.Errorf("invalid PCI slot format: %s", slot)
	}

	// The last part always contains "device.function"
	deviceFuncPart := parts[len(parts)-1]
	dfParts := strings.Split(deviceFuncPart, ".")
	if len(dfParts) != 2 {
		return SlotAddress{}, fmt.Errorf("invalid device.function format in %s", slot)
	}

	// Parse Device
	dev, err := strconv.ParseUint(dfParts[0], 16, 8)
	if err != nil {
		return SlotAddress{}, fmt.Errorf("invalid device value in %s: %v", slot, err)
	}
	device.Device = uint32(dev)

	// Parse Function
	funcVal, err := strconv.ParseUint(dfParts[1], 16, 8)
	if err != nil {
		return SlotAddress{}, fmt.Errorf("invalid function value in %s: %v", slot, err)
	}
	device.Function = uint32(funcVal)

	// Determine if domain is present
	if len(parts) == 3 {
		// Domain is present: [domain]:bus:device.function
		domain, err := strconv.ParseUint(parts[0], 16, 16)
		if err != nil {
			return SlotAddress{}, fmt.Errorf("invalid domain value in %s: %v", slot, err)
		}
		device.Domain = uint32(domain)

		bus, err := strconv.ParseUint(parts[1], 16, 8)
		if err != nil {
			return SlotAddress{}, fmt.Errorf("invalid bus value in %s: %v", slot, err)
		}
		device.Bus = uint32(bus)

	} else {
		// Domain is omitted (assumed 0000): bus:device.function
		device.Domain = 0 // Default to domain 0000
		bus, err := strconv.ParseUint(parts[0], 16, 8)
		if err != nil {
			return SlotAddress{}, fmt.Errorf("invalid bus value in %s: %v", slot, err)
		}
		device.Bus = uint32(bus)
	}

	return device, nil
}
