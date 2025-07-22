package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/olekukonko/tablewriter"
)

func listStacks(includeIncompatible bool) {
	stacksJson, err := snapctl.Get("stacks").Document().Run()
	if err != nil {
		log.Fatalf("Error loading stacks: %v\n", err)
	}

	stacks, err := parseStacksJson(stacksJson)
	if err != nil {
		log.Fatalf("Error parsing stacks: %v\n", err)
	}
	printStacks(stacks, includeIncompatible)
}

func parseStacksJson(stacksJson string) (map[string]types.ScoredStack, error) {
	var stacksOption map[string]map[string]types.ScoredStack
	err := json.Unmarshal([]byte(stacksJson), &stacksOption)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %v\n", err)
	}
	if stacks, ok := stacksOption["stacks"]; ok {
		return stacks, nil
	}
	return nil, fmt.Errorf("no stacks found")
}

func printStacks(stacks map[string]types.ScoredStack, includeIncompatible bool) error {

	var headers []string
	if includeIncompatible {
		headers = []string{"Stack Name", "Vendor", "Description", "Notes", "Compatible"}
	} else {
		headers = []string{"Stack Name", "Vendor", "Description", "Notes"}
	}
	data := [][]string{headers}

	// Iterate map in alphabetical order
	keys := slices.Collect(maps.Keys(stacks))
	slices.Sort(keys)

	for _, stackName := range keys {
		stack := stacks[stackName]
		stackInfo := []string{stack.Name, stack.Vendor, stack.Description, strings.Join(stack.Notes, ", ")}
		if includeIncompatible {
			if stack.Compatible {
				stackInfo = append(stackInfo, "Yes")
				data = append(data, stackInfo)
			} else {
				stackInfo = append(stackInfo, "No")
				data = append(data, stackInfo)
			}
		} else {
			if stack.Compatible {
				data = append(data, stackInfo)
			}
		}
	}

	if len(data) == 1 {
		if includeIncompatible {
			return fmt.Errorf("no stacks found")
		} else {
			return fmt.Errorf("no compatible stacks found")
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(data[0])
	err := table.Bulk(data[1:])
	if err != nil {
		return fmt.Errorf("error adding data to table: %v", err)
	}
	err = table.Render()
	if err != nil {
		return fmt.Errorf("error rendering table: %v", err)
	}
	return nil
}
