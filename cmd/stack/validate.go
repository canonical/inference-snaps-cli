package main

import (
	"fmt"

	"github.com/canonical/stack-utils/pkg/validate"
	"github.com/fatih/color"
)

func validateStackManifests(manifestFiles ...string) {
	for _, manifestPath := range manifestFiles {
		err := validate.Stack(manifestPath)
		if err != nil {
			fmt.Printf(color.RedString("x")+" %s: %s\n", manifestPath, err)
		} else {
			fmt.Printf(color.GreenString("✓")+" %s: valid\n", manifestPath)
		}
	}
}
