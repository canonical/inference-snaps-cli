package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canonical/stack-utils/pkg/selector"
	"github.com/canonical/stack-utils/pkg/types"

	"github.com/fatih/color"
)

type EngineSelection struct {
	Engines   []types.ScoredStack `json:"engines"`
	TopEngine string              `json:"top-engine"`
}

func main() {
	var enginesDir string
	var listAll bool
	var prettyOutput bool

	flag.StringVar(&enginesDir, "engines", "engines", "Override the path to the engines directory.")
	flag.BoolVar(&listAll, "all", false, "List all available engines.")
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
			log.Printf(color.RedString("x %s - %s"), engine.Name, strings.Join(engine.Notes, ", "))
		} else {
			log.Printf(color.GreenString("✓ %s - score = %d"), engine.Name, engine.Score)
		}
	}

	topEngine, err := selector.TopEngine(scoredEngines)
	if err != nil {
		log.Fatal(err)
	}
	engineSelection.TopEngine = topEngine.Name
	log.Printf(color.GreenString("Top engine: %s - score = %d"), topEngine.Name, topEngine.Score)

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
