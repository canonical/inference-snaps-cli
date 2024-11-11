package gpu

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGetLocalLsHw(t *testing.T) {
	lsHw, err := GetHostLsHw()
	if err != nil {
		t.Fatalf(err.Error())
	}

	log.Println(string(lsHw))
}

func TestParseLsHw(t *testing.T) {
	lsHw, err := GetHostLsHw()
	if err != nil {
		t.Fatalf(err.Error())
	}

	gpus, err := ParseLsHw(lsHw)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("%+v\n", gpus)
}

func TestParseLsHwDl380p(t *testing.T) {
	lsHw, err := os.ReadFile("test_data/hp-dl380p-gen8.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	gpuInfo, err := ParseLsHw(lsHw)
	if err != nil {
		t.Fatalf(err.Error())
	}

	jsonData, _ := json.MarshalIndent(gpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLsHwRpi5p(t *testing.T) {
	lsHw, err := os.ReadFile("test_data/rpi-5.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	gpuInfo, err := ParseLsHw(lsHw)
	if err != nil {
		t.Fatalf(err.Error())
	}

	jsonData, _ := json.MarshalIndent(gpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}

func TestParseLsHwXps13(t *testing.T) {
	lsHw, err := os.ReadFile("test_data/xps13-gen10.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	gpuInfo, err := ParseLsHw(lsHw)
	if err != nil {
		t.Fatalf(err.Error())
	}

	jsonData, _ := json.MarshalIndent(gpuInfo, "", "  ")
	fmt.Println(string(jsonData))
}
