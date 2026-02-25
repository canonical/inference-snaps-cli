package commands

import (
	"testing"
)

func TestVersionCleaning(t *testing.T) {
	if cleanVersionString("") != "unset" {
		t.Fatal("Expected empty version to be cleaned to 'unset'")
	}

	if cleanVersionString("1.2.3") != "1.2.3" {
		t.Fatal("Expected non-empty version to be unchanged")
	}
	if cleanVersionString("(devel)") != "(devel)" {
		t.Fatal("Expected '(devel)' version to be unchanged")
	}
	if cleanVersionString("1.0.0+dirty") != "1.0.0+dirty" {
		t.Fatal("Expected dirty version to be unchanged")
	}
}
