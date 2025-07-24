package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/canonical/go-snapctl/env"
)

var stacksDir = env.Snap() + "/stacks"

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalln("No subcommand provided, expected one of: select, load, download, validate")
	}

	if err := run(os.Args[2:]); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(args []string) error {

	switch args[1] {

	case "use":
		useCmd := flag.NewFlagSet("use", flag.ExitOnError)
		useAuto := useCmd.Bool("auto", false, "Automatically select a compatible stack")
		useAssumeYes := useCmd.Bool("assume-yes", false, "Assume yes for downloading new components")

		useCmd.Parse(os.Args[2:])

		if *useAuto {
			if len(useCmd.Args()) != 0 {
				fmt.Println("Error: cannot specify stack with --auto flag")
				os.Exit(1)
			}
			err := autoSelectStacks(*useAssumeYes)
			if err != nil {
				fmt.Println("Error: failed to automatically set used stack:", err)
				os.Exit(1)
			}
		} else {
			stack := useCmd.Args()
			if len(stack) == 1 {
				err := useStack(stack[0], *useAssumeYes)
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

	case "get":
		return get(args)

	case "set":
		return set(args)

	case "unset":
		return unset(args)

	case "load":
		if len(args) != 0 {
			return fmt.Errorf("'load' sub-command does not support arguments")
		}

		return load()

	case "download":
		if len(args) != 0 {
			return fmt.Errorf("'download' sub-command does not support arguments")
		}

		return download()

	// stack validate stacks/*/stack.yaml
	case "validate":
		return validate(args)

	case "list":
		listCmd.Parse(os.Args[2:])
		listStacks(*listAll)

	case "info":
		infoCmd.Parse(os.Args[2:])
		if len(infoCmd.Args()) < 1 {
			fmt.Println("Error: a stack name is required")
			os.Exit(1)
		}
		if len(infoCmd.Args()) != 1 {
			fmt.Println("Error: only one stack name can be specified")
			os.Exit(1)
		}
		stackInfo(infoCmd.Args()[0])

	default:
		return fmt.Errorf("unexpected sub-command: '%s'", args[1])
	}

	return nil

}
