package memory

import (
	"encoding/json"
	"testing"
)

func TestInfo(t *testing.T) {
	info, err := Info()
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	t.Log(string(jsonData))
}
