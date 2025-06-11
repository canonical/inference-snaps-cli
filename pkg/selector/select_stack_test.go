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
	//"../../test_data/hardware_info/i7-2600k.json", // Old CPU that does not have any of the necessary flags
	"../../test_data/hardware_info/mustang.json",
	"../../test_data/hardware_info/nuc11-i5-1145G7.json",
	"../../test_data/hardware_info/xps13-7390.json",
	"../../test_data/hardware_info/xps13-9350.json",
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
			scoredStacks, err := ScoreStacks(hardwareInfo, allStacks)
			if err != nil {
				t.Fatal(err)
			}
			topStack, err := TopStack(scoredStacks)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Found stack %s", topStack.Name)
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

	data, err = os.ReadFile("../../test_data/stacks/example-cpu/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var stack types.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	score, _, err := checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Matching score: %d", score)

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
	score, _, err = checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if score > 0 {
		t.Fatalf("Stack should not match if avx2 is not available. Score=%d", score)
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

	data, err = os.ReadFile("../../test_data/stacks/example-cpu-avx512/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var currentStack types.Stack
	err = yaml.Unmarshal(data, &currentStack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	score, reasons, err := checkStack(hardwareInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}
	if score == 0 {
		t.Fatalf("Stack should match if avx512 is available: %v", reasons)
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
	score, _, err = checkStack(hardwareInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}
	if score != 0 {
		t.Fatalf("Stack should not match if avx512 is not available. Score=%d", score)
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

func TestIntelDiscreteGpu(t *testing.T) {
	file, err := os.Open("../../test_data/hardware_info/i7-2600k+arc-a580.json")
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

	data, err = os.ReadFile("../../test_data/stacks/intel-gpu/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var currentStack types.Stack
	err = yaml.Unmarshal(data, &currentStack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	score, reasons, err := checkStack(hardwareInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}

	if score == 0 {
		t.Fatalf("Stack should match: %s", strings.Join(reasons, ",\n"))
	}
}

func TestMissingIntelDiscreteGpu(t *testing.T) {
	file, err := os.Open("../../test_data/hardware_info/xps13-7390.json")
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

	data, err = os.ReadFile("../../test_data/stacks/intel-gpu/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var currentStack types.Stack
	err = yaml.Unmarshal(data, &currentStack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	score, _, err := checkStack(hardwareInfo, currentStack)
	if err != nil {
		t.Fatal(err)
	}

	if score != 0 {
		t.Fatalf("Stack requiring GPU should not match device without GPU")
	}
}

func TestAmpere(t *testing.T) {
	var ampereHwInfos = []string{
		"../../test_data/hardware_info/ampere-one-x-mocked.json",
	}

	for _, hwInfoFile := range ampereHwInfos {
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

			data, err = os.ReadFile("../../test_data/stacks/ampere/stack.yaml")
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
				t.Fatalf("Ampere stack should match ampere hardware: %v", reasons)
			}
			t.Logf("Matching score: %d", score)

		})
	}
}

func TestAmpereNotMatchPi(t *testing.T) {
	file, err := os.Open("../../test_data/hardware_info/raspberry-pi-5.json")
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

	data, err = os.ReadFile("../../test_data/stacks/ampere/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var stack types.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	// Pi hardware should be invalid for ampere stack
	result, reasons, err := checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result != 0 {
		t.Fatal("Ampere stack should not match Raspberry Pi hardware")
	} else {
		t.Logf("Ampere stack does not match on Pi 5. Reason: %v", reasons)
	}
}

func TestNpuStackMatch(t *testing.T) {
	file, err := os.Open("../../test_data/hardware_info/xps13-9350.json")
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

	data, err = os.ReadFile("../../test_data/stacks/intel-npu/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var stack types.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	result, reasons, err := checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result == 0 {
		t.Fatalf("NPU stack should match npu hardware: %v", reasons)
	}
	t.Logf("Matching score: %d", result)
}
