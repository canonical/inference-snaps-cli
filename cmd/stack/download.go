package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
)

func downloadRequiredComponents() {
	// get stack snap option
	stackName, err := snapctl.Get("stack").Run()
	if err != nil {
		fmt.Println("Error getting stack from snap options:", err)
		os.Exit(1)
	}
	if stackName == "" {
		fmt.Println("Stack snap option is empty")
		os.Exit(1)
	}

	// get stacks.<new-stack> snap option for the list of components
	stackJson, err := snapctl.Get("stacks." + stackName).Run()
	if err != nil {
		fmt.Println("Error getting stack definition from snap options:", err)
		os.Exit(1)
	}
	var stack types.ScoredStack
	err = json.Unmarshal([]byte(stackJson), &stack)
	if err != nil {
		fmt.Println("Error deserializing stack definition from snap options:", err)
		os.Exit(1)
	}

	// install components
	for _, component := range stack.Components {
		fmt.Println("Installing component:", component)
		err = snapctl.InstallComponents(component).Run()
		if err != nil {
			if strings.Contains(err.Error(), "cannot install components for a snap that is unknown to the store") {
				fmt.Printf("Skipped: Install a local build: sudo snap install <path to %s component>\n", component)
			} else if strings.Contains(err.Error(), "already installed") {
				continue
			} else {
				fmt.Println("Error installing component:", err)
				os.Exit(1)
			}
		}
	}
}
