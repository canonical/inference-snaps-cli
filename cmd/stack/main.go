package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canonical/go-snapctl/env"
)

var stacksDir = env.Snap() + "/stacks"

func main() {
	// stack use [--yes] [--auto]
	// stack use [--yes] [<stack>]
	useCmd := flag.NewFlagSet("use", flag.ExitOnError)
	useAuto := useCmd.Bool("auto", false, "Automatically select a compatible stack")
	useYes := useCmd.Bool("yes", false, "Assume yes for downloading new components")

	// stack load
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)

	// stack download
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)

	// stack validate
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected a subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "use":
		useCmd.Parse(os.Args[2:])

		if *useAuto {
			if len(useCmd.Args()) != 0 {
				fmt.Println("Error: cannot specify stack with --auto flag")
				os.Exit(1)
			}
			err := autoSelectStacks(*useYes)
			if err != nil {
				fmt.Println("Error: failed to automatically set used stack:", err)
				os.Exit(1)
			}
		} else {
			stack := useCmd.Args()
			if len(stack) == 1 {
				err := useStack(stack[0], *useYes)
				if err != nil {
					fmt.Println("Error: failed use stack:", err)
					os.Exit(1)
				}
			} else if len(stack) == 0 {
				fmt.Println("Error: stack name not specified")
				os.Exit(1)
			} else {
				fmt.Println("Error: too many arguments")
				os.Exit(1)
			}
		}

	case "load":
		loadCmd.Parse(os.Args[2:])
		loadStacksToSnapOptions()

	case "download":
		downloadCmd.Parse(os.Args[2:])
		downloadRequiredComponents()

	// stack validate stacks/*/stack.yaml
	case "validate":
		validateCmd.Parse(os.Args[2:])
		stackFiles := validateCmd.Args()
		if len(stackFiles) == 0 {
			fmt.Println("Error: no stack manifest specified")
			os.Exit(1)
		}

		validateStackManifests(stackFiles...)

	default:
		fmt.Println("unexpected subcommands")
		os.Exit(1)
	}

}
