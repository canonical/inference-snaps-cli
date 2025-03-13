package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
	slog "github.com/canonical/go-snapctl/log"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info"
	"github.com/canonical/ml-snap-utils/pkg/selector"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

var stacksDir = env.Snap + "/stacks"

func main() {
	slog.SetComponentName("stack")

	// stack select [--auto]
	// stack select [<stack>]
	selectCmd := flag.NewFlagSet("select", flag.ExitOnError)
	selectAuto := selectCmd.Bool("auto", false, "Automatically select a compatible stack")

	// stack load
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)

	// stack download
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)

	if len(os.Args) < 2 {
		slog.Fatalf("expected a subcommands")
	}

	switch os.Args[1] {

	case "select":
		selectCmd.Parse(os.Args[2:])

		if *selectAuto {
			if len(selectCmd.Args()) != 0 {
				slog.Fatal("Error: cannot specify stack with --auto flag")
			}
			autoSelectStacks()
		} else {
			// for now falling here is an error, regardless of other arguments
			selectStack := selectCmd.Args()
			if len(selectStack) == 0 {
				slog.Fatal("Error: stack name not specified")
			} else {
				slog.Fatal("Not implemented! Use 'snap <snap> set stack=<stack>' instead")
			}
		}

	case "load":
		loadCmd.Parse(os.Args[2:])
		loadStacksToSnapOptions()

	case "download":
		downloadCmd.Parse(os.Args[2:])
		downloadRequiredComponents()

	default:
		fmt.Println("unexpected subcommands")
		os.Exit(1)
	}

}

func loadStacksToSnapOptions() {
	slog.Info("Loading stacks to snap options ...")

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		slog.Fatalf("Error loading stacks: %v", err)
	}

	// set all stacks as snap options
	// TODO: change to also handle stack deletions
	for _, stack := range allStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			slog.Fatalf("Error serializing stacks: %v", err)
		}

		err = snapctl.Set("stacks."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			slog.Fatalf("Error setting stacks option: %v", err)
		}
	}
}

func autoSelectStacks() {
	slog.Info("Automatically selecting a compatible stack ...")

	connected, err := snapctl.IsConnected("hardware-observe").Run()
	if err != nil {
		slog.Fatalf("Error checking hardware-observer connection: %v", err)
	}
	if !connected {
		slog.Info("Error: hardware-observe interface (https://snapcraft.io/docs/hardware-observe-interface) isn't connected.")
		slog.Info("This is required for hardware detection.")
		slog.Fatalf("Please connect and try again: sudo snap connect %s:hardware-observe", env.SnapName) // TODO: change to SnapInstanceName
	}

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		slog.Fatalf("Error loading stacks: %v", err)
	}

	slog.Infof("Found %d stacks", len(allStacks))

	// get hardware info
	hardwareInfo, err := hardware_info.Get(false)
	if err != nil {
		slog.Fatalf("Error getting hardware info: %v", err)
	}

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

	// set all scored stacks as snap options
	for _, stack := range scoredStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			slog.Fatalf("Error serializing stacks: %v", err)
		}

		err = snapctl.Set("stacks."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			slog.Fatalf("Error setting stacks option: %v", err)
		}
	}

	// find top stack
	topStack, err := selector.TopStack(scoredStacks)
	if err != nil {
		// FIXME: If no matching stacks are found, installation of this snap will fail
		slog.Fatal(err)
	}

	// set top stack name as snap option
	err = snapctl.Set("stack", topStack.Name).String().Run()
	if err != nil {
		slog.Fatalf("Error setting stack: %v", err)
	}

	// set snap options from stack configurations
	for confKey, confVal := range topStack.Configurations {
		valJson, err := json.Marshal(confVal)
		if err != nil {
			slog.Fatalf("Error serializing configuration %s: %v - %v", confKey, confVal, err)
		}
		err = snapctl.Set(confKey, string(valJson)).Document().Run()
		if err != nil {
			slog.Fatalf("can't set snap option: %v", err)
		}
	}

	slog.Infof("Selected stack for your hardware configuration: %s", topStack.Name)
}

func downloadRequiredComponents() {
	slog.Infof("Downloading required components ...")

	// get stack snap option
	stackName, err := snapctl.Get("stack").Run()
	if err != nil {
		slog.Fatalf("Error getting stack from snap options: %v", err)
	}
	if stackName == "" {
		slog.Fatal("Stack snap option is empty")
	}

	// get stacks.<new-stack> snap option for the list of components
	stackJson, err := snapctl.Get("stacks." + stackName).Run()
	if err != nil {
		slog.Fatalf("Error getting stack definition from snap options: %v", err)
	}
	var stack types.ScoredStack
	err = json.Unmarshal([]byte(stackJson), &stack)
	if err != nil {
		slog.Fatalf("Error deserializing stack definition from snap options: %v", err)
	}

	// install components
	for _, component := range stack.Components {
		err = snapctl.InstallComponents(component).Run()
		if err != nil {
			if strings.Contains(err.Error(), "cannot install components for a snap that is unknown to the store") {
				slog.Infof("Skip component installation. Install a local build: sudo snap install <path to %s component>", component)
			} else if strings.Contains(err.Error(), "already installed") {
				slog.Debugf("Skip component installation: already installed: %s", component)
			} else {
				slog.Fatalf("Error installing component: %v", err)
			}
		}
	}
}
