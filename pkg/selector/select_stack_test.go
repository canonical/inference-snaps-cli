package selector

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

var matchingHwStack = map[string]string{
	"../../test_data/hardware_info/amd-ryzen7-5700g.json": "../../test_data/stacks/example-cpu/stack.yaml",
	"../../test_data/hardware_info/amd-ryzen9-7900.json":  "../../test_data/stacks/example-cpu-avx512/stack.yaml",
	//"../../test_data/hardware_info/ampere-altra.json":        "../../test_data/stacks/ampere/stack.yaml", // need to add altra stack
	"../../test_data/hardware_info/ampere-one-x-mocked.json": "../../test_data/stacks/ampere/stack.yaml",
	"../../test_data/hardware_info/dell-r730xd.json":         "../../test_data/stacks/example-cpu/stack.yaml",
	"../../test_data/hardware_info/hp-dl380p-gen8.json":      "../../test_data/stacks/example-memory/stack.yaml",
	"../../test_data/hardware_info/i7-2600k+arc-a580.json":   "../../test_data/stacks/intel-gpu/stack.yaml",
	//"../../test_data/hardware_info/i7-2600k.json":            "../../test_data/stacks/example-cpu/stack.yaml", // too old cpu
	"../../test_data/hardware_info/mustang.json":         "../../test_data/stacks/intel-gpu/stack.yaml",
	"../../test_data/hardware_info/nuc11-i5-1145G7.json": "../../test_data/stacks/example-cpu-avx512/stack.yaml",
	//"../../test_data/hardware_info/raspberry-pi-5.json":      "../../test_data/stacks/example-cpu/stack.yaml", // no stacks for generic arm
	"../../test_data/hardware_info/xps13-7390.json": "../../test_data/stacks/example-cpu/stack.yaml",
	"../../test_data/hardware_info/xps13-9350.json": "../../test_data/stacks/intel-npu/stack.yaml",
}

func TestMatchingHwStack(t *testing.T) {
	for hwInfoFile, stackFile := range matchingHwStack {
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

			data, err = os.ReadFile(stackFile)
			if err != nil {
				t.Fatal(err)
			}

			var stack types.Stack
			err = yaml.Unmarshal(data, &stack)
			if err != nil {
				t.Fatal(err)
			}

			// Valid hardware for stack
			score, reasons, err := checkStack(hardwareInfo, stack)
			if err != nil {
				t.Fatal(err)
			}
			if score == 0 {
				t.Fatalf("Stack should match: %v", reasons)
			}
			t.Logf("Matching score: %d", score)
		})
	}
}

var nonMatchingHwStack = map[string]string{
	//"../../test_data/stacks/ampere/stack.yaml":             "../../test_data/hardware_info/raspberry-pi-5.json",
	//"../../test_data/stacks/example-cpu/stack.yaml":        "../../test_data/hardware_info/ampere-altra.json",
	//"../../test_data/stacks/example-cpu-avx512/stack.yaml": "../../test_data/hardware_info/i7-2600k+arc-a580.json",
	//"../../test_data/stacks/example-memory/stack.yaml":     "../../test_data/hardware_info/xps13-7390.json",
	//"../../test_data/stacks/generic-cuda/stack.yaml":       "../../test_data/hardware_info/raspberry-pi-5.json",
	//"../../test_data/stacks/intel-gpu/stack.yaml":          "../../test_data/hardware_info/raspberry-pi-5.json",
	"../../test_data/stacks/intel-npu/stack.yaml": "../../test_data/hardware_info/xps13-7390.json",
}

func TestStackNotMatchingHw(t *testing.T) {
	for stackFile, hwInfoFile := range nonMatchingHwStack {
		t.Run(stackFile, func(t *testing.T) {
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

			data, err = os.ReadFile(stackFile)
			if err != nil {
				t.Fatal(err)
			}

			var stack types.Stack
			err = yaml.Unmarshal(data, &stack)
			if err != nil {
				t.Fatal(err)
			}

			score, _, err := checkStack(hardwareInfo, stack)
			if err != nil {
				t.Fatal(err)
			}
			if score != 0 {
				t.Fatalf("Stack should not match: %s", hwInfoFile)
			}
			t.Logf("Matching score: %d", score)
		})
	}
}

func TestFindStackEmpty(t *testing.T) {
	hwInfo := types.HwInfo{
		Memory: &types.MemoryInfo{
			TotalRam:  200000000,
			TotalSwap: 200000000,
		},
		Disk: map[string]*types.DirStats{
			"/var/lib/snapd/snaps": &types.DirStats{
				Total: 0,
				Avail: 400000000,
			},
		},
	}

	allStacks, err := LoadStacksFromDir("../../test_data/stacks")
	if err != nil {
		t.Fatal(err)
	}
	scoredStacks, err := ScoreStacks(hwInfo, allStacks)
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
	hwInfo.Disk = make(map[string]*types.DirStats)
	hwInfo.Disk["/"] = &dirStat
	hwInfo.Disk["/var/lib/snapd/snaps"] = &dirStat

	stackDisk := "300M"
	stack := types.Stack{DiskSpace: &stackDisk}

	result, reasons, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result == 0 {
		t.Fatalf("disk should be enough: %v", reasons)
	}

	dirStat.Avail = 100000000
	result, reasons, err = checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result > 0 {
		t.Fatalf("disk should NOT be enough: %v", reasons)
	}
}

func TestMemoryCheck(t *testing.T) {
	hwInfo := types.HwInfo{
		Memory: &types.MemoryInfo{
			TotalRam:  200000000,
			TotalSwap: 200000000,
		},
	}

	stackMemory := "300M"
	stack := types.Stack{Memory: &stackMemory}

	result, reasons, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result == 0 {
		t.Fatalf("memory should be enough: %v", reasons)
	}

	hwInfo.Memory.TotalRam = 100000000
	result, reasons, err = checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result > 0 {
		t.Fatal("memory should NOT be enough")
	}
}

func TestNoCpuInHwInfo(t *testing.T) {
	hwInfo := types.HwInfo{
		// All fields are nil or zero
	}

	data, err := os.ReadFile("../../test_data/stacks/example-cpu-avx512/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var currentStack types.Stack
	err = yaml.Unmarshal(data, &currentStack)
	if err != nil {
		t.Fatal(err)
	}

	// No memory in hardware info
	_, _, err = checkStack(hwInfo, currentStack)
	if err == nil {
		t.Fatalf("No Memory in hardware_info should return err")
	}

	hwInfo.Memory = &types.MemoryInfo{
		TotalRam:  17000000000,
		TotalSwap: 2000000000,
	}

	// No disk space in hardware info
	_, _, err = checkStack(hwInfo, currentStack)
	if err == nil {
		t.Fatal("No Disk space in hardware_info should return err")
	}

	hwInfo.Disk = make(map[string]*types.DirStats)
	hwInfo.Disk["/"] = &types.DirStats{
		Avail: 6000000000,
	}

	// No CPU in hardware info
	_, _, err = checkStack(hwInfo, currentStack)
	if err == nil {
		t.Fatal("No CPU in hardware_info should return err")
	}
}
