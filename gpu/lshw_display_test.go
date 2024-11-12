package gpu

import (
	"encoding/json"
	"os"
	"testing"
)

func TestGetHostLsHw(t *testing.T) {
	lsHw, err := GetHostLsHw()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(string(lsHw))
}

func TestParseHostLsHw(t *testing.T) {
	lsHw, err := GetHostLsHw()
	if err != nil {
		t.Fatalf(err.Error())
	}

	gpus, err := ParseLsHw(lsHw)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Logf("%+v\n", gpus)
}

var testFiles = []string{
	"test_data/hp-dl380p-gen8.json",
	"test_data/rpi-5.json",
	"test_data/xps13-gen10.json",
}

func TestParseLsHw(t *testing.T) {
	for _, file := range testFiles {
		t.Run(file, func(t *testing.T) {
			lsHw, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf(err.Error())
			}

			gpuInfo, err := ParseLsHw(lsHw)
			if err != nil {
				t.Fatalf(err.Error())
			}

			jsonData, err := json.MarshalIndent(gpuInfo, "", "  ")
			if err != nil {
				t.Fatalf(err.Error())
			}

			t.Log(string(jsonData))
		})
	}
}
