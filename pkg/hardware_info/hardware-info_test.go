package hardware_info

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/canonical/stack-utils/pkg/hardware_info/cpu"
	"github.com/canonical/stack-utils/pkg/hardware_info/pci"
	"github.com/canonical/stack-utils/pkg/types"

	"github.com/go-test/deep"
)

var devices = []string{
	//"hp-proliant-rl300-gen11-altra",
	//"hp-proliant-rl300-gen11-altra-max",
	//"i7-2600k+arc-a580",
	//"raspberry-pi-5",
	//"raspberry-pi-5+hailo-8",
	//"xps13-7390",
}

func TestGetFromFiles(t *testing.T) {
	for _, device := range devices {
		t.Run(device, func(t *testing.T) {
			hwInfo, err := GetFromFiles(t, device, true)
			if err != nil {
				t.Error(err)
			}

			var hardwareInfo types.HwInfo
			devicePath := "../../test_data/devices/" + device + "/"
			hardwareInfoData, err := os.ReadFile(devicePath + "hardware-info.json")
			if err != nil {
				t.Fatal(err)
			}
			err = json.Unmarshal(hardwareInfoData, &hardwareInfo)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(hwInfo, hardwareInfo); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func GetFromFiles(t *testing.T, device string, friendlyNames bool) (types.HwInfo, error) {
	var hwInfo types.HwInfo

	devicePath := "../../test_data/devices/" + device + "/"

	// memory
	memoryData, err := os.ReadFile(devicePath + "memory.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(memoryData, &hwInfo.Memory)
	if err != nil {
		t.Fatal(err)
	}

	// disk
	diskData, err := os.ReadFile(devicePath + "disk.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(diskData, &hwInfo.Disk)
	if err != nil {
		t.Fatal(err)
	}

	// cpu
	unameMachine, err := os.ReadFile(devicePath + "uname-m.txt")
	if err != nil {
		t.Fatal(err)
	}
	procCpuInfo, err := os.ReadFile(devicePath + "cpuinfo.txt")
	if err != nil {
		t.Fatal(err)
	}
	cpuInfo, err := cpu.InfoFromData(string(procCpuInfo), string(unameMachine))
	if err != nil {
		t.Fatal(err)
	}
	hwInfo.Cpus = cpuInfo

	// pci
	pciData, err := os.ReadFile(devicePath + "lspci.txt")
	if err != nil {
		t.Fatal(err)
	}
	pciDevices, err := pci.DevicesFromData(string(pciData), friendlyNames)
	if err != nil {
		t.Fatal(err)
	}
	hwInfo.PciDevices = pciDevices

	// Additional properties - we append these directly from a file, as we can not run the vendor specific tools on the machine
	addPropsFile := devicePath + "additional-properties.json"
	_, err = os.Stat(addPropsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist. Skipping additional properties
		} else {
			t.Fatalf("error checking file '%s': %v\n", addPropsFile, err)
		}
	} else {
		var addProps map[string]map[string]string
		addPropsData, err := os.ReadFile(devicePath + "additional-properties.json")
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(addPropsData, &addProps)
		if err != nil {
			t.Fatal(err)
		}
		for i, pciDevice := range hwInfo.PciDevices {
			if val, ok := addProps[pciDevice.Slot]; ok {
				hwInfo.PciDevices[i].AdditionalProperties = val
			}
		}
	}

	return hwInfo, nil
}
