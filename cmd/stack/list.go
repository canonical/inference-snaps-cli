package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var (
	listAll bool
)

func init() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available stacks",
		// Long:  "",
		Args: cobra.NoArgs,
		RunE: list,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&listAll, "all", false, "include incompatible stacks")

	rootCmd.AddCommand(cmd)
}

func list(_ *cobra.Command, _ []string) error {
	return listStacks(listAll)
}

func listStacks(includeIncompatible bool) error {
	stacksJson, err := snapctl.Get("stacks").Document().Run()
	if err != nil {
		return fmt.Errorf("error loading stacks: %v", err)
	}

	stacks, err := parseStacksJson(stacksJson)
	if err != nil {
		return fmt.Errorf("error parsing stacks: %v", err)
	}

	err = printStacks(stacks, includeIncompatible)
	if err != nil {
		return fmt.Errorf("error printing list: %v", err)
	}

	return nil
}

func printStacks(stacks []types.ScoredStack, includeIncompatible bool) error {

	var headers []string
	if includeIncompatible {
		headers = []string{"Stack Name", "Vendor", "Description", "Status", "Notes"}
	} else {
		headers = []string{"Stack Name", "Vendor", "Description", "Status"}
	}
	data := [][]string{headers}

	// Sort by Score in descending order
	sort.Slice(stacks, func(i, j int) bool {
		return stacks[i].Score > stacks[j].Score
	})

	for _, stack := range stacks {
		stackInfo := []string{stack.Name, stack.Vendor, stack.Description}
		stackStatus := fmt.Sprintf("grade: %s\nscore: %d", stack.Grade, stack.Score)
		if includeIncompatible {
			if stack.Score == 0 {
				stackStatus = "incompatible\n" + stackStatus
			} else {
				stackStatus = "compatible\n" + stackStatus
			}
			stackInfo = append(stackInfo, stackStatus, strings.Join(stack.Notes, ", "))
			data = append(data, stackInfo)
		} else {
			stackInfo = append(stackInfo, stackStatus)
			if stack.Score > 0 {
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

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
		})),
		tablewriter.WithMaxWidth(80),
	)
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
