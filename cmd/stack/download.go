package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
)

const (
	snapdUnknownSnapError = "cannot install components for a snap that is unknown to the store"
	snapdTimeoutError     = "timeout exceeded while waiting for response"
)

func download() error {
	return downloadRequiredComponents()
}

func downloadRequiredComponents() error {
	// get stack snap option
	stackName, err := snapctl.Get("stack").Run()
	if err != nil {
		return fmt.Errorf("error getting 'stack' from snap options: %s", err)
	}
	if stackName == "" {
		return fmt.Errorf("'stack' snap option is empty")
	}

	// get stacks.<new-stack> snap option for the list of components
	stackJson, err := snapctl.Get("stacks." + stackName).Run()
	if err != nil {
		return fmt.Errorf("error getting 'stacks' from snap options: %s", err)
	}
	var stack types.ScoredStack
	err = json.Unmarshal([]byte(stackJson), &stack)
	if err != nil {
		return fmt.Errorf("error deserializing 'stacks': %s", err)
	}

	downloadComponents(stack.Components)
}

func downloadComponents(components []string) {
	// install components
	// Messages presented to the user should use the term "download" for snapctl install +component.
	for _, component := range components {
		stopProgress := startProgressDots("Downloading " + component + " ")
		err := snapctl.InstallComponents(component).Run()
		stopProgress()
		if err != nil {
			if strings.Contains(err.Error(), snapdUnknownSnapError) {
				fmt.Printf("Error: snap not known to the store. Install a local build of component: %s", component)
				continue
			} else if strings.Contains(err.Error(), snapdTimeoutError) {
				msg := "timeout exceeded while waiting for download of: " + component +
					"\nPlease monitor the progress using the 'snap changes' command and continue when the component installation is complete."
				return fmt.Errorf(msg)
			} else if strings.Contains(err.Error(), "already installed") {
				continue
			} else {
				return fmt.Errorf("error downloading component: %s: %s", component, err)
			}
		}
		fmt.Println("Downloaded " + component)
	}

	return nil
}

func startProgressDots(prefix string) (stop func()) {
	dots := []string{".", "..", "..."}
	s := spinner.New(dots, time.Second)
	s.Prefix = prefix
	s.Start()

	return s.Stop
}
