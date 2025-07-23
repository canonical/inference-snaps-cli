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

	case "select": // deprecated - do not touch
		selectCmd := flag.NewFlagSet("select", flag.ExitOnError)
		selectAuto := selectCmd.Bool("auto", false, "Automatically select a compatible stack")
		selectCmd.Parse(os.Args[2:])

		if *selectAuto {
			if len(selectCmd.Args()) != 0 {
				fmt.Println("Error: cannot specify stack with --auto flag")
				os.Exit(1)
			}
			autoSelectStacks()
		} else {
			stack := selectCmd.Args()
			if len(stack) == 1 {
				selectStack(stack[0])
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

	default:
		return fmt.Errorf("unexpected sub-command: '%s'", args[1])
	}

	return nil

}
