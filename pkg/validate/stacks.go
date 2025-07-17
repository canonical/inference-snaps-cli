package validate

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/canonical/stack-utils/pkg/types"
	"github.com/canonical/stack-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

func Stack(manifestFilePath string) error {
	_, err := os.Stat(manifestFilePath)
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	yamlData, err := os.ReadFile(manifestFilePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return validateStackYaml(manifestFilePath, yamlData)
}

func validateStackYaml(manifestFilePath string, yamlData []byte) error {

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

	return validateStackStruct(manifestFilePath, stack)
}

func validateStackStruct(manifestFilePath string, stack types.Stack) error {
	if stack.Name == "" {
		return fmt.Errorf("required field is not set: name")
	}

	// The name inside the manifest file should be the same as the one in the file path
	expectedPath := filepath.Join(stack.Name, "stack.yaml")
	if !strings.HasSuffix(manifestFilePath, expectedPath) {
		return fmt.Errorf("stack dir name should equal name in manifest: %s != %s", manifestFilePath, stack.Name)
	}

	if stack.Description == "" {
		return fmt.Errorf("required field is not set: description")
	}

	if stack.Vendor == "" {
		return fmt.Errorf("required field is not set: vendor")
	}

	if stack.Grade == "" {
		return fmt.Errorf("required field is not set: grade")
	}
	if stack.Grade != "stable" && stack.Grade != "devel" {
		return fmt.Errorf("grade should be 'stable' or 'devel'")
	}

	if stack.Memory != nil {
		_, err := utils.StringToBytes(*stack.Memory)
		if err != nil {
			return fmt.Errorf("error parsing memory: %v", err)
		}
	}

	if stack.DiskSpace != nil {
		_, err := utils.StringToBytes(*stack.DiskSpace)
		if err != nil {
			return fmt.Errorf("error parsing disk space: %v", err)
		}
	}

	if _, ok := stack.Configurations["engine"]; !ok {
		return fmt.Errorf("required field is not set: config.engine")
	}

	if _, ok := stack.Configurations["model"]; !ok {
		return fmt.Errorf("required field is not set: config.model")
	}

	return stackDevices(stack.Devices)
}
