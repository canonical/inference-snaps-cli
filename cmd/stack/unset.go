package main

import (
	"fmt"

	"github.com/canonical/go-snapctl"
)

func unset(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected one config key as input")
	}

	return unsetValue(args[0])
}

func unsetValue(key string) error {
	err := snapctl.Unset(key).Run()
	if err != nil {
		return fmt.Errorf("error unsetting value of '%s': %v", key, err)
	}

	return nil
}
