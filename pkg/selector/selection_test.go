package selector

import (
	"fmt"
	"os"
	"testing"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

// Test that the expected stack is chosen from a list of stacks
// device, stacks, expected stack

type machineTopStack struct {
	machine  string
	stacks   []string
	topStack string
}

var topStackSets = []machineTopStack{
	{
		machine:  "ampere-one-m-banshee-12",
		stacks:   []string{"ampere", "arm-neon"},
		topStack: "ampere",
	},
	{
		machine:  "hp-proliant-rl300-gen11-altra",
		stacks:   []string{"ampere-altra", "arm-neon"},
		topStack: "ampere-altra",
	},
	{
		machine:  "mustang",
		stacks:   []string{"cpu-avx1", "cpu-avx2", "cpu-avx512", "intel-cpu", "intel-gpu"},
		topStack: "intel-gpu",
	},
	{
		machine:  "xps13-7390",
		stacks:   []string{"cpu-avx1", "cpu-avx2"},
		topStack: "cpu-avx2",
	},
	{
		machine:  "xps13-7390",
		stacks:   []string{"cpu-avx1", "cpu-avx2", "intel-cpu"},
		topStack: "intel-cpu",
	},
}

func TestTopStack(t *testing.T) {
	for _, testSet := range topStackSets {
		t.Run(testSet.machine+"/"+testSet.topStack, func(t *testing.T) {
			var stacks []types.Stack
			for _, stackName := range testSet.stacks {
				stackManifestFile := fmt.Sprintf("../../test_data/stacks/%s/stack.yaml", stackName)
				data, err := os.ReadFile(stackManifestFile)
				if err != nil {
					t.Fatal(err)
				}

				var stack types.Stack
				err = yaml.Unmarshal(data, &stack)
				if err != nil {
					t.Fatal(err)
				}

				stacks = append(stacks, stack)
			}

			hardwareInfo, err := hardware_info.GetFromRawData(t, testSet.machine, true)
			if err != nil {
				t.Fatal(err)
			}

			scoredStacks, err := ScoreStacks(hardwareInfo, stacks)
			if err != nil {
				t.Fatal(err)
			}

			topStack, err := TopStack(scoredStacks)
			if err != nil {
				t.Fatal(err)
			}

			if topStack.Name != testSet.topStack {
				t.Errorf("Top stack name: %s, expected: %s", topStack.Name, testSet.topStack)
			}
		})
	}
}
