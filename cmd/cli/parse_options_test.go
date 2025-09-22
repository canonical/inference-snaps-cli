package main

import (
	"os"
	"testing"
)

func TestParseEnginesJson(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/engines.json")
	if err != nil {
		t.Fatal(err)
	}

	_, err = parseEnginesJson(string(data))
	if err != nil {
		t.Fatal(err)
	}
}
