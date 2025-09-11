package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canonical/stack-utils/pkg/engines"
	"github.com/canonical/stack-utils/pkg/selector"
	"github.com/canonical/stack-utils/pkg/types"

	"github.com/fatih/color"
)

type EngineSelection struct {
	Engines   []engines.ScoredManifest `json:"engines"`
	TopEngine string                   `json:"top-engine"`
}

func main() {
	var enginesDir string
	var prettyOutput bool

	flag.StringVar(&enginesDir, "engines", "engines", "Override the path to the engines directory.")
	flag.BoolVar(&prettyOutput, "pretty", false, "Pretty print JSON.")

	flag.Parse()

	// Read json piped in from the hardware-info app
	var hardwareInfo types.HwInfo

	err := json.NewDecoder(os.Stdin).Decode(&hardwareInfo)
	if err != nil {
		log.Fatal(err)
	}

	allEngines, err := selector.LoadManifestsFromDir(enginesDir)
	if err != nil {
		log.Fatal(err)
	}
	scoredEngines, err := selector.ScoreEngines(hardwareInfo, allEngines)
	if err != nil {
		log.Fatal(err)
	}

	var engineSelection EngineSelection

	// Print summary on STDERR
	for _, engine := range scoredEngines {
		engineSelection.Engines = append(engineSelection.Engines, engine)

		if engine.Score == 0 {
			fmt.Fprintf(os.Stderr, "❌ %s - not compatible: %s\n", engine.Name, strings.Join(engine.Notes, ", "))
		} else if engine.Grade != "stable" {
			fmt.Fprintf(os.Stderr, "🟠 %s - score = %d, grade = %s\n", engine.Name, engine.Score, engine.Grade)
		} else {
			fmt.Fprintf(os.Stderr, "✅ %s - compatible, score = %d\n", engine.Name, engine.Score)
		}
	}

	selectedEngine, err := selector.TopEngine(scoredEngines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding top engine: %v", err)
	}
	engineSelection.TopEngine = selectedEngine.Name

	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Fprintf(os.Stderr, greenBold("Selected engine for your hardware configuration: %s\n\n"), selectedEngine.Name)

	var resultStr []byte
	if prettyOutput {
		resultStr, err = json.MarshalIndent(engineSelection, "", "  ")
	} else {
		resultStr, err = json.Marshal(engineSelection)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", resultStr)
}
