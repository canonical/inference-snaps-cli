package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

func info(args []string) error {
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)
	infoCmd.Parse(args)

	if len(args) != 1 {
		return fmt.Errorf("expected one stack name as input")
	}

	stackName := args[0]
	// check: can this arg be present and empty?
	// if stackName == "" {
	// 	return fmt.Errorf("stack name cannot be empty")
	// }

	stackInfo(stackName)

	return nil
}

func stackInfo(stackName string) {
	stackJson, err := snapctl.Get("stacks." + stackName).Document().Run()
	if err != nil {
		log.Fatalf("Error loading stack: %v\n", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		log.Fatalf("Error parsing stack: %v\n", err)
	}
	err = printStackInfo(stack)
	if err != nil {
		log.Fatalf("Error printing stack info: %v\n", err)
	}
}

func printStackInfo(stack types.ScoredStack) error {
	stackYaml, err := yaml.Marshal(stack)
	if err != nil {
		return fmt.Errorf("error converting stack to yaml: %v", err)
	}

	err = quick.Highlight(os.Stdout, string(stackYaml), "yaml", "terminal", "colorful")
	if err != nil {
		return fmt.Errorf("error formatting yaml: %v", err)
	}

	return nil
}
