package memory

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetInfo(t *testing.T) {
	info, err := GetInfo()
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%+v", info)

	jsonData, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(jsonData))
}
