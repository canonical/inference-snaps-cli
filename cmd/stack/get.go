package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/canonical/go-snapctl"
)

func get(args []string) error {
	flag := flag.NewFlagSet("get", flag.ExitOnError)
	help := helpFlag(flag)

	fmt.Println("get:", args)

	flag.Parse(args)
	fmt.Println("get args:", flag.Args())

	// keep the non-flag arguments
	args = flag.Args()

	flag.Usage = func() {
		log.Printf("Usage:\n\t%s get <key>", snapInstanceName)
	}

	if *help {
		flag.Usage()
		return nil
	}

	if len(args) != 1 {
		return fmt.Errorf("expected one config key as input, got %q", strings.Join(args, " "))
	}

	return getValue(args[0])
}

func getValue(key string) error {
	value, err := snapctl.Get(key).Run()
	if err != nil {
		return fmt.Errorf("error getting value of '%s': %v", key, err)
	}

	if value == "" {
		return fmt.Errorf("no value set for key '%s'", key)
	}

	// print config value
	fmt.Println(value)

	return nil
}
