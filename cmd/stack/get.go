package main

import (
	"fmt"

	"github.com/canonical/go-snapctl"
)

func get(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected one config key as input")
	}

	return getValue(args[0])
}

func getValue(key string) error {
	value, err := snapctl.Get(key).Run()
	if err != nil {
		return fmt.Errorf("error getting value of '%s': %v", key, err)
	}

	if value == "" {
		return fmt.Errorf("no value set for key '%s'", key)
	}

	// print config value
	fmt.Println(value)

	return nil
}
