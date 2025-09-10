package validate

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/stack-utils/pkg/types"
	"github.com/canonical/stack-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

func Engine(manifestFilePath string) error {

	if !strings.HasSuffix(manifestFilePath, "engine.yaml") {
		return fmt.Errorf("stack manifest file must be called engine.yaml: %s", manifestFilePath)
	}

	_, err := os.Stat(manifestFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stack manifest file does not exist: %s", manifestFilePath)
		}
		return fmt.Errorf("error getting file info: %v", err)
	}

	yamlData, err := os.ReadFile(manifestFilePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Get stack name from path
	stackName := engineNameFromPath(manifestFilePath)

	return validateEngineManifestYaml(stackName, yamlData)
}

func engineNameFromPath(manifestFilePath string) string {
	parts := utils.SplitPathIntoDirectories(manifestFilePath)
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2] // second last part: stack-name/engine.yaml
}

func validateEngineManifestYaml(expectedStackName string, yamlData []byte) error {
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

	return validateEngineStruct(expectedStackName, stack)
}

func validateEngineStruct(expectedEngineName string, stack types.Stack) error {
	if stack.Name == "" {
		return fmt.Errorf("required field is not set: name")
	}

	// Only do engine name matching test if expected name is set
	if expectedEngineName != "" {
		if stack.Name != expectedEngineName {
			return fmt.Errorf("engine directory name should match name in manifest: %s != %s", expectedEngineName, stack.Name)
		}
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

	for key, val := range stack.Configurations {
		if !utils.IsPrimitive(val) {
			return fmt.Errorf("configuration field %s is not a primitive value: %v", key, val)
		}
	}

	return stackDevices(stack.Devices)
}
