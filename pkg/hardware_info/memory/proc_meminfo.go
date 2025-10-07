package memory

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/canonical/famous-models-cli/pkg/types"
)

func hostProcMemInfo() (string, error) {
	// cat /proc/meminfo
	memInfoBytes, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", fmt.Errorf("error reading /proc/meminfo: %v", err)
	}
	return string(memInfoBytes), nil
}

func parseProcMemInfo(memInfoString string) (types.MemoryInfo, error) {
	var memInfo = types.MemoryInfo{}

	lines := strings.Split(memInfoString, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		fields := strings.SplitN(line, ":", 2)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0]) // remove \t between key and colon
		value := strings.TrimSpace(fields[1])

		switch key {
		case "MemTotal":
			valueBytes, err := procStringToBytes(value)
			if err != nil {
				return memInfo, fmt.Errorf("error parsing MemTotal: %v", err)
			}
			memInfo.TotalRam = uint64(valueBytes)
		case "SwapTotal":
			valueBytes, err := procStringToBytes(value)
			if err != nil {
				return memInfo, fmt.Errorf("error parsing SwapTotal: %v", err)
			}
			memInfo.TotalSwap = uint64(valueBytes)
		}
	}
	return memInfo, nil
}

func procStringToBytes(s string) (int64, error) {
	s = strings.TrimSpace(s)

	if strings.HasSuffix(s, "kB") {
		s = strings.TrimSuffix(s, "kB")
		s = strings.TrimSpace(s)
		kbValue, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing kB value: %v", err)
		}
		return kbValue * 1024, nil
	} else {
		bValue, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing byte value: %v", err)
		}
		return bValue, nil
	}
}
