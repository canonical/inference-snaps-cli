package cpu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ProcCpuInfo contains general information about a system CPU found in /proc/cpuinfo.
type ProcCpuInfo struct {
	Processor    int64    // %d - kernel defines it as long long
	ModelName    *string  // %s
	BogoMips     float64  // %lu.%02lu
	Features     []string // space separated strings
	Implementer  uint64   // 0x%02x
	Architecture uint64   // constant int
	Variant      uint64   // 0x%x
	Part         uint64   // 0x%03x
	Revision     uint64   // %d
}

func procCpuInfo() ([]ProcCpuInfo, error) {
	procCpuInfoBytes, err := hostProcCpuInfoBytes()
	if err != nil {
		return nil, err
	}

	return parseProcCpuInfo(procCpuInfoBytes)
}

func hostProcCpuInfoBytes() ([]byte, error) {
	// cat /proc/cpuinfo
	cpuInfoBytes, err := os.ReadFile("/proc/cpuinfo")
	return cpuInfoBytes, err
}

func parseProcCpuInfo(cpuInfoBytes []byte) ([]ProcCpuInfo, error) {
	cpuInfoString := string(cpuInfoBytes)

	// amd64 machines have a static string set
	// see https://github.com/torvalds/linux/blob/master/arch/arm64/kernel/cpuinfo.c#L267
	if strings.Contains(cpuInfoString, "CPU architecture: 8\n") {
		return parseProcCpuInfoAmd64(cpuInfoString)
	}

	// TODO x86 see https://github.com/torvalds/linux/blob/ec7714e4947909190ffb3041a03311a975350fe0/arch/x86/kernel/cpu/proc.c#L63

	return nil, fmt.Errorf("cpu info parsing not implemented for this architecture")
}

func parseProcCpuInfoAmd64(cpuInfoString string) ([]ProcCpuInfo, error) {
	var parsedCpus []ProcCpuInfo

	lines := strings.Split(cpuInfoString, "\n")
	previousLineEmpty := false
	var parsedCpu ProcCpuInfo

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			previousLineEmpty = true
			continue
		}

		// New cpu block
		if previousLineEmpty {
			previousLineEmpty = false
			parsedCpus = append(parsedCpus, parsedCpu)
			parsedCpu = ProcCpuInfo{}
		}

		fields := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(fields[0]) // remove \t between key and colon
		value := strings.TrimSpace(fields[1])

		switch key {
		// Formatting strings above cases are from torvalds/linux/blob/master/arch/arm64/kernel/cpuinfo.c

		// "processor\t: %d\n"
		case "processor":
			processorIndex, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			parsedCpu.Processor = processorIndex

		// "model name\t: ARMv8 Processor rev %d (%s)\n"
		case "model name":
			modelName := strings.TrimSpace(value)
			parsedCpu.ModelName = &modelName

		// BogoMIPS\t: %lu.%02lu\n
		case "BogoMIPS":
			bogoMips, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			parsedCpu.BogoMips = bogoMips

		// "Features\t:"+" %s"
		case "Features":
			flags := strings.Split(value, " ")
			parsedCpu.Features = append(parsedCpu.Features, flags...)

		// "CPU implementer\t: 0x%02x\n"
		case "CPU implementer":
			implementer, err := strconv.ParseUint(value, 0, 8) // use base 0 to allow parser to detect and remove 0x prefix
			if err != nil {
				return nil, err
			}
			parsedCpu.Implementer = implementer

		// "CPU architecture: 8\n"
		case "CPU architecture":
			architecture, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, err
			}
			parsedCpu.Architecture = architecture

		// "CPU variant\t: 0x%x\n"
		case "CPU variant":
			variant, err := strconv.ParseUint(value, 0, 64)
			if err != nil {
				return nil, err
			}
			parsedCpu.Variant = variant

		// "CPU part\t: 0x%03x\n"
		case "CPU part":
			part, err := strconv.ParseUint(value, 0, 16)
			if err != nil {
				return nil, err
			}
			parsedCpu.Part = part

		// "CPU revision\t: %d\n\n"
		case "CPU revision":
			revision, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, err
			}
			parsedCpu.Revision = revision
		}
	}
	return parsedCpus, nil
}
