package validate

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/canonical/stack-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

func Stack(stackFile string) error {
	_, err := os.Stat(stackFile)
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	yamlData, err := os.ReadFile(stackFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return validateStackYaml(yamlData)
}

func validateStackYaml(yamlData []byte) error {

	yamlData = bytes.TrimSpace(yamlData)
	if len(yamlData) == 0 {
		return errors.New("empty yaml data")
	}

	var stack types.Stack

	yamlDecoder := yaml.NewDecoder(bytes.NewReader(yamlData))

	// Error if there are unknown fields in the yaml
	yamlDecoder.KnownFields(true)

	// We depend on the yaml unmarshal to check field types
	if err := yamlDecoder.Decode(&stack); err != nil {
		return fmt.Errorf("error decoding: %v", err)
	}

	return nil
}
