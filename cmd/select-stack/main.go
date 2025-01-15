package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/canonical/ml-snap-utils/pkg/selector"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func main() {
	var stacksDir string
	var listAll bool
	var prettyOutput bool
	var summary bool

	flag.StringVar(&stacksDir, "stacks", "stacks", "Override the path to the stacks directory.")
	flag.BoolVar(&listAll, "all", false, "List all available stacks.")
	flag.BoolVar(&prettyOutput, "pretty", false, "Pretty print JSON.")
	flag.BoolVar(&summary, "summary", false, "Show summary of stacks.")

	flag.Parse()

	// Read json piped in from the hardware-info app
	var hardwareInfo types.HwInfo

	err := json.NewDecoder(os.Stdin).Decode(&hardwareInfo)
	if err != nil {
		log.Fatal(err)
	}

	var resultStr []byte

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		log.Fatal(err)
	}
	scoredStacks, err := selector.ScoreStacks(hardwareInfo, allStacks)
	if err != nil {
		log.Fatal(err)
	}
	bestStack, err := selector.BestStack(scoredStacks)
	if err != nil {
		log.Fatal(err)
	}

	var result interface{}
	if listAll {
		result = scoredStacks
	} else {
		result = bestStack
	}

	if summary {
		for _, stack := range scoredStacks {
			if stack.Score == 0 {
				log.Printf("Stack %s not selected: %s", stack.Name, stack.Comment)
			} else {
				log.Printf("Stack %s matches. Score = %d", stack.Name, stack.Score)
			}
		}
		log.Printf("Best stack: %s. Score = %d", bestStack.Name, bestStack.Score)

	} else {
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
}
