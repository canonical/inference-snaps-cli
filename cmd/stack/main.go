package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info"
	"github.com/canonical/ml-snap-utils/pkg/selector"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

var stacksDir = env.Snap + "/stacks"

func main() {
	// stack select [--auto]
	// stack select [<stack>]
	selectCmd := flag.NewFlagSet("select", flag.ExitOnError)
	selectAuto := selectCmd.Bool("auto", false, "Automatically select a compatible stack")

	// stack load
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)

	// stack download
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected a subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "select":
		selectCmd.Parse(os.Args[2:])

		if *selectAuto {
			if len(selectCmd.Args()) != 0 {
				fmt.Println("Error: cannot specify stack with --auto flag")
				os.Exit(1)
			}
			autoSelectStacks()
		} else {
			// for now falling here is an error, regardless of other arguments
			selectStack := selectCmd.Args()
			if len(selectStack) == 0 {
				fmt.Println("Error: stack name not specified")
				os.Exit(1)
			} else {
				fmt.Println("Not implemented! Use 'snap <snap> set stack=<stack>' instead")
				os.Exit(1)
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
	fmt.Println("Loading stacks to snap options ...")

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		fmt.Printf("Error loading stacks: %v", err)
		os.Exit(1)
	}

	// set all stacks as snap options
	// TODO: change to also handle stack deletions
	for _, stack := range allStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			fmt.Println("Error serializing stacks:", err)
			os.Exit(1)
		}

		err = snapctl.Set("stacks."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			fmt.Println("Error setting stacks option:", err)
			os.Exit(1)
		}
	}
}

func autoSelectStacks() {
	fmt.Println("Automatically selecting a compatible stack ...")

	connected, err := snapctl.IsConnected("hardware-observe").Run()
	if err != nil {
		fmt.Println("Error checking hardware-observer connection:", err)
		os.Exit(1)
	}
	if !connected {
		fmt.Println("Error: hardware-observe interface (https://snapcraft.io/docs/hardware-observe-interface) isn't connected.")
		fmt.Println("This is required for hardware detection.")
		fmt.Printf("Please connect and try again: sudo snap connect %s:hardware-observe\n", env.SnapInst)
		os.Exit(1)
	}

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		fmt.Println("Error loading stacks:", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d stacks\n", len(allStacks))

	// get hardware info
	hardwareInfo, err := hardware_info.Get(false)
	if err != nil {
		fmt.Println("Error getting hardware info:", err)
		os.Exit(1)
	}

	// score stacks
	scoredStacks, err := selector.ScoreStacks(hardwareInfo, allStacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, stack := range scoredStacks {
		if stack.Score == 0 {
			fmt.Printf("Stack %s not selected: %s\n", stack.Name, strings.Join(stack.Notes, ", "))
		} else {
			fmt.Printf("Stack %s matches. Score = %d\n", stack.Name, stack.Score)
		}
	}

	// set all scored stacks as snap options
	for _, stack := range scoredStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			fmt.Println("Error serializing stacks:", err)
			os.Exit(1)
		}

		err = snapctl.Set("stacks."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			fmt.Println("Error setting stacks option:", err)
			os.Exit(1)
		}
	}

	// find top stack
	topStack, err := selector.TopStack(scoredStacks)
	if err != nil {
		fmt.Println("Error selecting a stack:", err)
		os.Exit(1)
	}

	// set top stack name as snap option
	err = snapctl.Set("stack", topStack.Name).String().Run()
	if err != nil {
		fmt.Println("Error setting stack:", err)
		os.Exit(1)
	}

	// set snap options from stack configurations
	for confKey, confVal := range topStack.Configurations {
		valJson, err := json.Marshal(confVal)
		if err != nil {
			fmt.Printf("Error serializing configuration %s: %v - %v\n", confKey, confVal, err)
			os.Exit(1)
		}
		err = snapctl.Set(confKey, string(valJson)).Document().Run()
		if err != nil {
			fmt.Println("Error setting snap option:", err)
			os.Exit(1)
		}
	}

	fmt.Println("Selected stack for your hardware configuration:", topStack.Name)
}

func downloadRequiredComponents() {
	fmt.Println("Downloading required components ...")

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
		err = snapctl.InstallComponents(component).Run()
		if err != nil {
			if strings.Contains(err.Error(), "cannot install components for a snap that is unknown to the store") {
				fmt.Printf("Skip component installation. Install a local build: sudo snap install <path to %s component>\n", component)
			} else if strings.Contains(err.Error(), "already installed") {
				fmt.Println("Skip component installation: already installed:", component)
			} else {
				fmt.Println("Error installing component:", err)
				os.Exit(1)
			}
		}
	}
}
