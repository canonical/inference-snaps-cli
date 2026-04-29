package commands

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func GetFromRawDataWrapper(t *testing.T) (*types.HwInfo, error) {
	info, err := hardware_info.GetFromRawData(t, "dummy-machine", true, "../../../test_data")
	if err != nil {
		return nil, err
	}
	return info, nil
}
func Example_showMachineCommand_printMachineInfoJson() {
	cmd := showMachineCommand{format: "json"}
	testingT := &testing.T{}
	info, err := GetFromRawDataWrapper(testingT)
	if err != nil {
		panic(err)
	}

	if err := cmd.printMachineInfoJson(info); err != nil {
		panic(err)
	}

	// Output:
	// {
	//   "cpus": [
	//     {
	//       "architecture": "amd64",
	//       "manufacturer-id": "GenuineIntel",
	//       "flags": [
	//         "fpu",
	//         "vme",
	//         "de"
	//       ]
	//     }
	//   ],
	//   "memory": {
	//     "total-ram": 67012501504,
	//     "total-swap": 0
	//   },
	//   "disk": {
	//     "/var/lib/snapd/snaps": {
	//       "total": 1006451294208,
	//       "avail": 943543738368
	//     }
	//   },
	//   "pci": [
	//     {
	//       "slot": "0000:00:00.0",
	//       "bus-number": "0x0",
	//       "device-class": "0x600",
	//       "programming-interface": 0,
	//       "vendor-id": "0x8086",
	//       "device-id": "0x4637",
	//       "subvendor-id": "0x103C",
	//       "subdevice-id": "0x89C6",
	//       "vendor-name": "Intel Corporation",
	//       "subvendor-name": "Hewlett-Packard Company"
	//     }
	//   ]
	// }

}

func Example_showMachineCommand_printMachineInfoYaml() {
	cmd := showMachineCommand{format: "yaml"}
	testingT := &testing.T{}
	info, err := GetFromRawDataWrapper(testingT)
	if err != nil {
		panic(err)
	}

	if err := cmd.printMachineInfoYaml(info); err != nil {
		panic(err)
	}

	// Output:
	// cpus:
	//     - architecture: amd64
	//       manufacturer-id: GenuineIntel
	//       flags:
	//         - fpu
	//         - vme
	//         - de
	// memory:
	//     total-ram: 67012501504
	//     total-swap: 0
	// disk:
	//     /var/lib/snapd/snaps:
	//         total: 1006451294208
	//         avail: 943543738368
	// pci:
	//     - slot: "0000:00:00.0"
	//       bus-number: "0x0"
	//       device-class: "0x600"
	//       programming-interface: 0
	//       vendor-id: "0x8086"
	//       device-id: "0x4637"
	//       subvendor-id: "0x103C"
	//       subdevice-id: "0x89C6"
	//       vendor-name: Intel Corporation
	//       subvendor-name: Hewlett-Packard Company
}

func Test_printMachineInfo_unknownFormat(t *testing.T) {
	cmd := showMachineCommand{format: "xml"}
	info := &types.HwInfo{}

	err := cmd.printMachineInfo(info)
	if err == nil || err.Error() != `unknown format "xml"` {
		t.Errorf("expected error 'unknown format \"xml\"', got %v", err)
	}
}
