package runtimes

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"go.yaml.in/yaml/v4"
)

func Validate(manifestFilePath string) error {

	if !strings.HasSuffix(manifestFilePath, ManifestFilename) {
		return fmt.Errorf("manifest file must be called %s: %s", ManifestFilename, manifestFilePath)
	}

	_, err := os.Stat(manifestFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("manifest file does not exist: %s", manifestFilePath)
		}
		return fmt.Errorf("getting file info: %v", err)
	}

	yamlData, err := os.ReadFile(manifestFilePath)
	if err != nil {
		return fmt.Errorf("reading file: %v", err)
	}

	// Get runtime name from path
	runtimeName := runtimeNameFromPath(manifestFilePath)

	return validateManifestYaml(runtimeName, yamlData)
}

func runtimeNameFromPath(manifestFilePath string) string {
	parts := utils.SplitPathIntoDirectories(manifestFilePath)
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2] // second last part: runtime-name/runtime.yaml
}

func validateManifestYaml(expectedName string, yamlData []byte) error {
	yamlData = bytes.TrimSpace(yamlData)
	if len(yamlData) == 0 {
		return errors.New("empty yaml data")
	}

	var manifest Manifest

	yamlDecoder := yaml.NewDecoder(bytes.NewReader(yamlData))

	// Error if there are unknown fields in the yaml
	yamlDecoder.KnownFields(true)

	// We depend on the yaml unmarshal to check field types
	if err := yamlDecoder.Decode(&manifest); err != nil {
		return fmt.Errorf("decoding manifest: %v", err)
	}

	return manifest.validate(expectedName)
}

func (manifest Manifest) validate(expectedRuntimeName string) error {
	if manifest.Name == "" {
		return fmt.Errorf("required field is not set: name")
	}

	// Only do runtime name matching test if expected name is set
	if expectedRuntimeName != "" {
		if manifest.Name != expectedRuntimeName {
			return fmt.Errorf("runtime directory name should match name in manifest: %s != %s", expectedRuntimeName, manifest.Name)
		}
	}

	if len(manifest.Servers) == 0 {
		return fmt.Errorf("required field is not set: servers")
	}

	for name, server := range manifest.Servers {
		if err := server.validate(name); err != nil {
			return err
		}
	}

	if len(manifest.Environment) == 0 {
		return fmt.Errorf("required field is not set: environment")
	}

	if len(manifest.Components) == 0 {
		return fmt.Errorf("required field is not set: components")
	}

	return nil
}

func (server Server) validate(name string) error {
	if server.Protocol == "" {
		return fmt.Errorf("required field is not set for server %s: protocol", name)
	}

	if server.BasePath == "" {
		return fmt.Errorf("required field is not set for server %s: base-path", name)
	}

	return nil
}
