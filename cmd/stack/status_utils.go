package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
)

const (
	openAi = "openai"
)

type Status struct {
	Variant   string             `json:"variant" yaml:"variant"`
	Status    string             `json:"status" yaml:"status"`
	Endpoints map[string]url.URL `json:"endpoints" yaml:"endpoints"`
}

func scoredStacksFromOptions() ([]types.ScoredStack, error) {
	stacksJson, err := snapctl.Get("stacks").Document().Run()
	if err != nil {
		return nil, fmt.Errorf("error loading variants: %v", err)
	}

	stacksMap, err := parseStacksJson(stacksJson)
	if err != nil {
		return nil, fmt.Errorf("error parsing variants: %v", err)
	}

	// map to slice
	var stacks []types.ScoredStack
	for _, stack := range stacksMap {
		stacks = append(stacks, stack)
	}

	return stacks, nil
}

func selectedStackFromOptions() (types.ScoredStack, error) {
	selectedStackName, err := snapctl.Get("stack").Run()
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error loading selected variant: %v", err)
	}

	stackJson, err := snapctl.Get("stacks." + selectedStackName).Document().Run()
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error loading variant: %v", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error parsing variant: %v", err)
	}

	return stack, nil
}

func statusStruct() (*Status, error) {
	var statusStr Status

	// Find the selected stack
	stack, err := selectedStackFromOptions()
	if err != nil {
		return nil, fmt.Errorf("error loading selected variant: %v", err)
	}
	statusStr.Variant = stack.Name

	ssc, err := serverStatusCode(stack.Name)
	if err != nil {
		return nil, fmt.Errorf("error getting server status: %v", err)
	}
	switch ssc {
	case 0:
		statusStr.Status = "online"
	case 1:
		statusStr.Status = "starting"
	case 2:
		statusStr.Status = "offline"
	default:
		statusStr.Status = "unknown"
	}

	endpoints, err := serverApiUrls(stack)
	if err != nil {
		return nil, fmt.Errorf("error getting server api endpoints: %v", err)
	}
	statusStr.Endpoints = endpoints

	return &statusStr, nil
}

func serverApiUrls(stack types.ScoredStack) (map[string]url.URL, error) {
	// Build API URL
	apiBasePath := "v1"
	if val, ok := stack.Configurations["http.base-path"]; ok {
		apiBasePath, ok = val.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for base path: %v", val)
		}

	}
	httpPort, err := snapctl.Get("http.port").Run()
	if err != nil {
		return nil, fmt.Errorf("error getting http port: %v", err)
	}

	openaiHost := fmt.Sprintf("localhost:%s", httpPort)
	openaiUrl := url.URL{Scheme: "http", Host: openaiHost, Path: apiBasePath}
	return map[string]url.URL{openAi: openaiUrl}, nil

	// TODO add additional api endpoints like openvino on http://localhost:8080/v1
}

func serverStatusCode(stackName string) (int, error) {
	// Depend on existing check server scripts for status
	checkScript := os.ExpandEnv("$SNAP/stacks/" + stackName + "/check-server")
	cmd := exec.Command(checkScript)
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("error checking server: %v", err)
	}

	checkExitCode := 0
	if err := cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			checkExitCode = exitError.ExitCode()
		}
	}
	return checkExitCode, nil
}
