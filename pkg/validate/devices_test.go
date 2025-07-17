package validate

import (
	"testing"

	"github.com/canonical/stack-utils/pkg/types"
)

func TestStackDeviceType(t *testing.T) {
	device := types.StackDevice{}

	t.Run("Type CPU", func(t *testing.T) {
		device.Type = "cpu"
		err := stackDevice(device)
		if err != nil {
			t.Fatalf("Type cpu should be valid: %v", err)
		}
	})
	t.Run("Type GPU", func(t *testing.T) {
		device.Type = "gpu"
		err := stackDevice(device)
		if err != nil {
			t.Fatalf("Type gpu should be valid: %v", err)
		}
	})
	t.Run("Type NPU", func(t *testing.T) {
		device.Type = "npu"
		err := stackDevice(device)
		if err != nil {
			t.Fatalf("Type npu should be valid: %v", err)
		}
	})
	t.Run("Type empty", func(t *testing.T) {
		device.Type = ""
		err := stackDevice(device)
		if err != nil {
			t.Fatalf("Empty type should be valid: %v", err)
		}
	})
	t.Run("Type invalid", func(t *testing.T) {
		device.Type = "test"
		err := stackDevice(device)
		if err == nil {
			t.Fatalf("Invalid type should be invalid: %v", err)
		}
		t.Log(err)
	})
}
