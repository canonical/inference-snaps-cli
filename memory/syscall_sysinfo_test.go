package memory

import (
	"encoding/json"
	"testing"
)

func TestGetInfo(t *testing.T) {
	info, err := GetInfo()
	if err != nil {
		t.Fatalf(err.Error())
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(string(jsonData))
}
