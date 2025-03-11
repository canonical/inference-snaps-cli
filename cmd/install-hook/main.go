package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	slog "github.com/canonical/go-snapctl/log"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info"
	"github.com/canonical/ml-snap-utils/pkg/selector"
)

func main() {
	// get all stacks from "$SNAP"/stacks/*/stack.yaml and append to a single json array
	snapDir, ok := os.LookupEnv("SNAP")
	if !ok {
		slog.Fatalf("SNAP environment variable not set")
	}

	stacksDir := snapDir + "/stacks"
	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		slog.Fatalf("Error loading stacks: %v", err)
	}

	slog.Info("Found %d stacks", len(allStacks))

	stacksJson, err := json.Marshal(allStacks)
	if err != nil {
		slog.Fatalf("Error serializing stacks: %v", err)
	}

	// set all stacks as snap options under `stacks`
	err = snapctl.Set("stacks", string(stacksJson)).Document().Run()
	if err != nil {
		slog.Fatalf("Error setting stacks option: %v", err)
	}

	// get hardware info
	hardwareInfo, err := hardware_info.Get(false)

	// score stacks
	scoredStacks, err := selector.ScoreStacks(hardwareInfo, allStacks)
	if err != nil {
		slog.Fatal(err)
	}

	for _, stack := range scoredStacks {
		if stack.Score == 0 {
			slog.Infof("Stack %s not selected: %s", stack.Name, strings.Join(stack.Notes, ", "))
		} else {
			slog.Infof("Stack %s matches. Score = %d", stack.Name, stack.Score)
		}
	}

	// find top stack
	topStack, err := selector.TopStack(scoredStacks)
	if err != nil {
		slog.Fatal(err)
	}

	// set top stack name in `stack` snap option
	err = snapctl.Set("stack", topStack.Name).String().Run()
	if err != nil {
		slog.Fatalf("Error setting stack: %v", err)
	}

	// install components
	snapctl.InstallComponents(topStack.Components...)

	// get top stack configurations

	// set snap options for configurations
	configurationsString, err := json.Marshal(topStack.Configurations)
	if err != nil {
		slog.Fatalf("can't convert configurations to string: %v", err)
	}
	err = snapctl.Set("", string(configurationsString)).Document().Run()
	if err != nil {
		slog.Fatalf("can't set snap option: %v", err)
	}

	// set generic configurations
	err = snapctl.Set("http.port", "8080").Run()
	if err != nil {
		slog.Fatalf("can't set snap http.port: %v", err)
	}

	err = snapctl.Set("http.host", "127.0.0.1").Run()
	if err != nil {
		slog.Fatalf("can't set snap http.host: %v", err)
	}
}
