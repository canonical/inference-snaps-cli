package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/canonical/go-snapctl/env"
)

var (
	stacksDir        = env.Snap() + "/stacks"
	snapInstanceName = env.SnapInstanceName()
)

func main() {
	// disable logging timestamps
	log.SetFlags(0)

	if snapInstanceName == "" {
		snapInstanceName = "<model>"
	}

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalln("No subcommand provided, expected one of: select, load, download, validate")
	}

	if err := subcommand(flag.Args()); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func subcommand(args []string) error {
	fmt.Println("subcommand", args)

	switch args[0] {

	case "use":
		return use(args[1:])

	case "get":
		return get(args[1:])

	case "set":
		return set(args[1:])

	case "unset":
		return unset(args[1:])

	case "load":
		// todo: move to load
		if len(args) != 0 {
			return fmt.Errorf("'load' sub-command does not support arguments")
		}

		return load()

	// stack validate stacks/*/stack.yaml
	case "validate":
		return validate(args[1:])

	case "list":
		return list(args[1:])

	case "info":
		return info(args[1:])

	default:
		return fmt.Errorf("unexpected sub-command: %q", args[0])
	}
}

// helpFlag adds the help flag to a flagset
func helpFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("help", false, "Show usage")
}
