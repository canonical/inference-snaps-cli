package main

import (
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
		return use(args)

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

	// stack validate stacks/*/stack.yaml
	case "validate":
		return validate(args)

	case "list":
		return list(args)

	case "info":
		return info(args)

	default:
		return fmt.Errorf("unexpected sub-command: '%s'", args[1])
	}

	return nil

}
