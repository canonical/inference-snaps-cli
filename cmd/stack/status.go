package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/spf13/cobra"
)

var (
	statusAll bool
)

func init() {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the status of the model snap",
		// Long:  "",
		Args: cobra.NoArgs,
		RunE: status,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&statusAll, "all", false, "include hardware information")

	rootCmd.AddCommand(cmd)
}

func status(_ *cobra.Command, _ []string) error {
	return snapStatus(statusAll)
}

func snapStatus(showHardware bool) error {
	/*
		$ mymodel status
		Stack: intel-gpu (auto)
		Model: mymodel-7b-distill-int4
		Engine: openvino-model-server

		Server:
		    Status: starting|online|offline
		    OpenAI endpoint: http://localhost:8080/v1

		[--all]:
		Hardware and resources:
	*/
	scoredStacks, err := scoredStacksFromOptions()
	if err != nil {
		return fmt.Errorf("error loading scored stacks: %v", err)
	}
	autoStack, err := topStack(scoredStacks)
	if err != nil {
		return fmt.Errorf("error loading top stack: %v", err)
	}
	stack, err := selectedStackFromOptions()
	if err != nil {
		return fmt.Errorf("error loading selected stack: %v", err)
	}

	printStack(stack, stack.Name == autoStack.Name)
	fmt.Println("")
	printServer(stack)
	if showHardware {
		fmt.Println("")
		printHardware(stack)
	}
	fmt.Println("")

	return nil
}

func scoredStacksFromOptions() ([]types.ScoredStack, error) {
	stacksJson, err := snapctl.Get("stacks").Document().Run()
	if err != nil {
		return nil, fmt.Errorf("error loading stacks: %v", err)
	}

	stacksMap, err := parseStacksJson(stacksJson)
	if err != nil {
		return nil, fmt.Errorf("error parsing stacks: %v", err)
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
		return types.ScoredStack{}, fmt.Errorf("error loading selected stack: %v", err)
	}

	stackJson, err := snapctl.Get("stacks." + selectedStackName).Document().Run()
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error loading stack: %v", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error parsing stack: %v", err)
	}

	return stack, nil
}

func printStack(stack types.ScoredStack, auto bool) {
	autoString := ""
	if auto {
		autoString = " (auto)"
	}
	fmt.Printf("Stack: %s%s\n", stack.Name, autoString)

	if val, ok := stack.Configurations["model"]; ok {
		fmt.Printf("  Model: %s\n", val)
	}
	if val, ok := stack.Configurations["engine"]; ok {
		fmt.Printf("  Engine: %s\n", val)
	}
	if val, ok := stack.Configurations["multimodel-projector"]; ok {
		fmt.Printf("  Multimodal projector: %s\n", val)
	}
}

func printServer(stack types.ScoredStack) {
	apiBasePath := "v1"
	if val, ok := stack.Configurations["http.base-path"]; ok {
		apiBasePath = val.(string)
	}
	httpPort, err := snapctl.Get("http.port").Run()
	if err != nil {
		fmt.Printf("error loading http port: %v", err)
		return
	}

	// Depend on existing check server scripts for status
	checkExitCode := 0
	checkScript := os.ExpandEnv("$SNAP/stacks/" + stack.Name + "/check-server")
	cmd := exec.Command(checkScript)
	if err := cmd.Start(); err != nil {
		fmt.Printf("cmd.Start: %v", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			checkExitCode = exiterr.ExitCode()
		} else {
			fmt.Errorf("cmd.Wait: %v", err)
			return
		}
	}

	statusText := "online"
	switch checkExitCode {
	case 0:
		statusText = "online"
	case 1:
		statusText = "starting"
	case 2:
		statusText = "offline"
	default:
		statusText = fmt.Sprintf("unknown (exit code %d)", checkExitCode)
	}

	fmt.Printf("Server:\n")
	fmt.Printf("  Status: %s\n", statusText)
	fmt.Printf("  OpenAI endpoint: http://localhost:%s/%s\n", httpPort, apiBasePath)
}

func printHardware(stack types.ScoredStack) {

}
