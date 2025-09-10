package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/stack-utils/pkg/types"
)

func parseEnginesJson(enginesJson string) ([]types.ScoredStack, error) {
	var enginesOption map[string]map[string]types.ScoredStack
	err := json.Unmarshal([]byte(enginesJson), &enginesOption)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %v", err)
	}
	if enginesMap, ok := enginesOption["engines"]; ok {
		var enginesSlice []types.ScoredStack
		for _, engine := range enginesMap {
			enginesSlice = append(enginesSlice, engine)
		}
		return enginesSlice, nil
	}
	return nil, fmt.Errorf("no engines found")
}

func parseEngineJson(engineJson string) (types.ScoredStack, error) {
	var engineOption map[string]types.ScoredStack

	err := json.Unmarshal([]byte(engineJson), &engineOption)
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error parsing json: %v", err)
	}

	if len(engineOption) == 0 {
		return types.ScoredStack{}, fmt.Errorf("engine not found")
	}

	if len(engineOption) > 1 {
		return types.ScoredStack{}, fmt.Errorf("only one engine expected in json")
	}

	for _, engine := range engineOption {
		return engine, nil
	}

	return types.ScoredStack{}, fmt.Errorf("unexpected error occurred")
}
