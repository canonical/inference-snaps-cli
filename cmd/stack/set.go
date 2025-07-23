package main

import (
	"fmt"
	"strings"

	"github.com/canonical/go-snapctl"
)

func set(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected one key=value pair as input")
	}
	return setValue(args[0])
}

func setValue(keyValue string) error {
	if keyValue[0] == '=' {
		return fmt.Errorf("key must not start with an equal sign")
	}

	// The value itself can contain an equal sign, so we split only on the first occurrence
	parts := strings.SplitN(keyValue, "=", 2)
	key, value := parts[0], parts[1]

	err := snapctl.Set(key, value).Run()
	if err != nil {
		return fmt.Errorf("error setting value '%s' for '%s': %v", value, key, err)
	}

	return nil
}
