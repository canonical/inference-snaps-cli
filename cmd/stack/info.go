package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/canonical/stack-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

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

func parseStackJson(stackJson string) (types.ScoredStack, error) {
	var stackOption map[string]types.ScoredStack

	err := json.Unmarshal([]byte(stackJson), &stackOption)
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error parsing json: %v", err)
	}

	if len(stackOption) == 0 {
		return types.ScoredStack{}, fmt.Errorf("stack not found")
	}

	if len(stackOption) > 1 {
		return types.ScoredStack{}, fmt.Errorf("only one stack expected in json")
	}

	for _, stack := range stackOption {
		return stack, nil
	}

	return types.ScoredStack{}, fmt.Errorf("unexpected error occurred")
}

func printStackInfo(stack types.ScoredStack) error {
	fmt.Println("# General")
	fmt.Printf("Stack name:\t%s\n", stack.Name)
	fmt.Printf("Description:\t%s\n", stack.Description)
	fmt.Printf("Vendor:\t\t%s\n", stack.Vendor)
	fmt.Printf("Grade:\t\t%s\n", stack.Grade)
	fmt.Println()

	fmt.Println("# Compatibility")
	fmt.Printf("Compatible:\t%v\n", stack.Compatible)
	fmt.Printf("Score:\t\t%d\n", stack.Score)
	if len(stack.Notes) > 0 {
		fmt.Printf("Notes:\n")
		for _, note := range stack.Notes {
			fmt.Printf("- %s\n", note)
		}
	}
	fmt.Println()

	fmt.Println("# Components")
	for _, component := range stack.Components {
		fmt.Printf("- %s\n", component)
	}
	fmt.Println()

	fmt.Println("# Configurations")
	for key, value := range stack.Configurations {
		fmt.Printf("- %s: %v\n", key, value)
	}
	fmt.Println()

	fmt.Println("# Requirements")
	if stack.Memory != nil {
		memoryBytes, err := utils.StringToBytes(*stack.Memory)
		if err != nil {
			return fmt.Errorf("error parsing stack memory: %v", err)
		}
		fmt.Printf("Memory:\t\t%s\n", utils.FmtGigabytes(memoryBytes))
	}
	if stack.DiskSpace != nil {
		diskBytes, err := utils.StringToBytes(*stack.DiskSpace)
		if err != nil {
			return fmt.Errorf("error parsing stack disk space: %v", err)
		}
		fmt.Printf("Disk space:\t%s\n", utils.FmtGigabytes(diskBytes))
	}

	if len(stack.Devices.All) > 0 {
		fmt.Printf("All these devices are required:\n")
		devicesYaml, err := yaml.Marshal(stack.Devices.All)
		if err != nil {
			return fmt.Errorf("unable to format devices all: %v", err)
		}
		devicesYaml = bytes.TrimSpace(devicesYaml) // remove newline added by yaml marshaller
		fmt.Printf("%s\n", string(devicesYaml))
	}

	if len(stack.Devices.Any) > 0 {
		fmt.Printf("One or more of these devices is required:\n")
		devicesYaml, err := yaml.Marshal(stack.Devices.Any)
		if err != nil {
			return fmt.Errorf("unable to format devices any: %v", err)
		}
		devicesYaml = bytes.TrimSpace(devicesYaml)
		fmt.Printf("%s\n", string(devicesYaml))
	}

	return nil
}
