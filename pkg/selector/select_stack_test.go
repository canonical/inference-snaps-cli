package selector

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

var hwInfoFiles = []string{
	"../../test_data/hardware_info/amd-ryzen7-5700g.json",
	"../../test_data/hardware_info/amd-ryzen9-7900.json",
	"../../test_data/hardware_info/dell-r730xd.json",
	"../../test_data/hardware_info/hp-dl380p-gen8.json",
	"../../test_data/hardware_info/nuc11-i5-1145G7.json",
	"../../test_data/hardware_info/xps13-gen10.json",
}

func TestFindStack(t *testing.T) {
	for _, hwInfoFile := range hwInfoFiles {
		t.Run(hwInfoFile, func(t *testing.T) {
			file, err := os.Open(hwInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			data, err := io.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			var hardwareInfo types.HwInfo
			err = json.Unmarshal(data, &hardwareInfo)
			if err != nil {
				t.Fatal(err)
			}

			allStacks, err := LoadStacksFromDir("../../test_data/stacks")
			if err != nil {
				t.Fatal(err)
			}
			scoredStacks, err := FilterStacks(hardwareInfo, allStacks)
			if err != nil {
				t.Fatal(err)
			}
			topStack, err := TopStack(scoredStacks)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Found stack %s with size %d", topStack.Name, topStack.Size)
		})
	}
}

func TestFindStackEmpty(t *testing.T) {
	hwInfo := types.HwInfo{}

	allStacks, err := LoadStacksFromDir("../../test_data/stacks-empty")
	if err != nil {
		t.Fatal(err)
	}
	scoredStacks, err := FilterStacks(hwInfo, allStacks)
	if err != nil {
		t.Fatal(err)
	}
	topStack, err := TopStack(scoredStacks)
	if err == nil {
		t.Fatal("Empty stack dir should return an error for top stack")
	}
	if topStack != nil {
		t.Fatal("No stack should be found in empty stacks dir")
	}
}

func TestDiskCheck(t *testing.T) {
	dirStat := types.DirStats{
		Total: 0,
		Avail: 400000000,
	}
	hwInfo := types.HwInfo{}
	hwInfo.Memory = &types.MemoryInfo{TotalRam: 2000000}
	hwInfo.Disk = make(map[string]*types.DirStats)
	hwInfo.Disk["/"] = &dirStat
	hwInfo.Disk["/var/lib/snapd/snaps"] = &dirStat

	stack := types.Stack{DiskSpace: "300M", Memory: "1M"}

	compatible, notes, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if !compatible {
		t.Fatalf("disk should be enough: %s", strings.Join(notes, ", "))
	}

	dirStat.Avail = 100000000
	compatible, _, err = checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatal("Not enough disk should not be compatible")
	}
}

func TestMemoryCheck(t *testing.T) {
	hwInfo := types.HwInfo{}
	hwInfo.Memory = &types.MemoryInfo{TotalRam: 200000000, TotalSwap: 200000000}
	hwInfo.Disk = make(map[string]*types.DirStats)
	hwInfo.Disk["/var/lib/snapd/snaps"] = &types.DirStats{
		Total: 0,
		Avail: 400000000,
	}

	stack := types.Stack{Memory: "300M", DiskSpace: "1M"}

	compatible, notes, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if !compatible {
		t.Fatalf("memory should be enough: %s", strings.Join(notes, ", "))
	}

	hwInfo.Memory.TotalRam = 100000000
	compatible, _, err = checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatal("memory should NOT be enough")
	}
}

func TestCpuFlagsAvx2(t *testing.T) {
	file, err := os.Open("../../test_data/hardware_info/amd-ryzen7-5700g.json")
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	var hardwareInfo types.HwInfo
	err = json.Unmarshal(data, &hardwareInfo)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile("../../test_data/stacks/llamacpp-avx2/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var stack types.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	compatible, notes, err := checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if !compatible {
		t.Fatalf("cpu should be compatible: %s", strings.Join(notes, ", "))
	}

	file, err = os.Open("../../test_data/hardware_info/hp-dl380p-gen8.json")
	if err != nil {
		t.Fatal(err)
	}

	data, err = io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &hardwareInfo)
	if err != nil {
		t.Fatal(err)
	}

	// Invalid hardware for stack
	compatible, notes, err = checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatalf("Stack should not match if avx2 is not available")
	}
}

func TestCpuFlagsAvx512(t *testing.T) {
	file, err := os.Open("../../test_data/hardware_info/amd-ryzen9-7900.json")
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	var hardwareInfo types.HwInfo
	err = json.Unmarshal(data, &hardwareInfo)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile("../../test_data/stacks/llamacpp-avx512/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var currentStack types.Stack
	err = yaml.Unmarshal(data, &currentStack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	compatible, _, err := checkStack(hardwareInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}
	if !compatible {
		t.Fatal("cpu should be compatible")
	}

	file, err = os.Open("../../test_data/hardware_info/hp-dl380p-gen8.json")
	if err != nil {
		t.Fatal(err)
	}

	data, err = io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &hardwareInfo)
	if err != nil {
		t.Fatal(err)
	}

	// Invalid hardware for stack
	compatible, _, err = checkStack(hardwareInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatal("Stack should not match if avx512 is not available")
	}
}

func TestNoCpuInHwInfo(t *testing.T) {
	hwInfo := types.HwInfo{
		// All fields are nil or zero
	}

	data, err := os.ReadFile("../../test_data/stacks/llamacpp-avx512/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var currentStack types.Stack
	err = yaml.Unmarshal(data, &currentStack)
	if err != nil {
		t.Fatal(err)
	}

	// No memory in hardware info
	compatible, _, err := checkStack(hwInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatal("No Memory in hardware_info should not be compatible")
	}
	//t.Log(err)

	hwInfo.Memory = &types.MemoryInfo{
		TotalRam:  17000000000,
		TotalSwap: 2000000000,
	}

	// No disk space in hardware info
	compatible, _, err = checkStack(hwInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatal("No Disk space in hardware_info should not be compatible")
	}

	hwInfo.Disk = make(map[string]*types.DirStats)
	hwInfo.Disk["/"] = &types.DirStats{
		Avail: 6000000000,
	}

	// No CPU in hardware info
	compatible, _, err = checkStack(hwInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatal("No CPU in hardware_info should not be compatible")
	}
}
