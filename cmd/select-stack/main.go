package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/selector"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func main() {
	var stacksDir string
	var listAll bool
	var prettyOutput bool

	flag.StringVar(&stacksDir, "stacks", "stacks", "Override the path to the stacks directory.")
	flag.BoolVar(&listAll, "all", false, "List all available stacks.")
	flag.BoolVar(&prettyOutput, "pretty", false, "Pretty print JSON.")

	flag.Parse()

	// Read json piped in from the hardware-info app
	var hardwareInfo types.HwInfo

	err := json.NewDecoder(os.Stdin).Decode(&hardwareInfo)
	if err != nil {
		log.Fatal(err)
	}

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		log.Fatal(err)
	}
	checkedStacks, err := selector.CheckStacks(allStacks, hardwareInfo)
	if err != nil {
		log.Fatal(err)
	}

	var result types.StackSelection

	// Append compatible stacks to result. Print summary on STDERR
	for _, stack := range checkedStacks {
		if stack.Compatible {
			result.Stacks = append(result.Stacks, stack)
			log.Printf("%s is compatible. Size: %d", stack.Name, stack.Size)
		} else {
			log.Printf("%s not compatible. Reasons: %s", stack.Name, strings.Join(stack.Notes, ", "))
		}
	}

	topStack, err := selector.TopStack(checkedStacks)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Top stack: %s. Size = %d", topStack.Name, topStack.Size)
	result.TopStack = topStack.Name

	var resultStr []byte
	if prettyOutput {
		resultStr, err = json.MarshalIndent(result, "", "  ")
	} else {
		resultStr, err = json.Marshal(result)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", resultStr)
}
