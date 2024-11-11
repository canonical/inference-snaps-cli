package gpu

import (
	"encoding/json"
	"os/exec"
)

func GetHostLsHw() ([]byte, error) {
	out, err := exec.Command("lshw", "-json", "-C", "display").Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ParseLsHw(input []byte) ([]Display, error) {
	displays := make([]Display, 0)

	var lsHwParsed []map[string]interface{}
	err := json.Unmarshal(input, &lsHwParsed)
	if err != nil {
		return displays, err
	}

	for _, lsHwDisplay := range lsHwParsed {
		display := Display{}

		// We use the Vendor and Product IDs reported by `lshw`
		// What will happen if a vendor introduces a new ID/Name?
		// Vendor and product IDs are from https://pci-ids.ucw.cz/
		if vendor, ok := lsHwDisplay["vendor"]; ok {
			display.Vendor = vendor.(string)
		}
		if product, ok := lsHwDisplay["product"]; ok {
			display.Product = product.(string)
		}
		displays = append(displays, display)
	}

	return displays, nil
}
