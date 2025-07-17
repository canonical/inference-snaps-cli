package validate

import (
	"fmt"

	"github.com/canonical/stack-utils/pkg/types"
)

func stackDevices(devices types.StackDevices) error {
	for i, device := range devices.All {
		err := stackDevice(device)
		if err != nil {
			return fmt.Errorf("invalid device: all %d/%d: %v", i+1, len(devices.All), err)
		}
	}

	for i, device := range devices.Any {
		err := stackDevice(device)
		if err != nil {
			return fmt.Errorf("invalid device: any %d/%d: %v", i+1, len(devices.All), err)
		}
	}

	return nil
}

func stackDevice(device types.StackDevice) error {
	switch device.Type {
	case "cpu":
		return cpu(device)
	case "gpu":
		return gpu(device)
	case "npu":
		return npu(device)
	case "":
		return typelessDevice(device)
	default:
		return fmt.Errorf("invalid device type: %v", device.Type)
	}
}

func gpu(device types.StackDevice) error {
	return nil
}

func npu(device types.StackDevice) error {
	return nil
}

func typelessDevice(device types.StackDevice) error {
	return nil
}
