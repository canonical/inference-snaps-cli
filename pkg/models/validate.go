package models

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

	// Get model name from path
	modelId := modelIdFromPath(manifestFilePath)

	return validateManifestYaml(modelId, yamlData)
}

func modelIdFromPath(manifestFilePath string) string {
	parts := utils.SplitPathIntoDirectories(manifestFilePath)
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2] // second last part: model-id/model.yaml
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

func (manifest Manifest) validate(expectedModelId string) error {
	if manifest.ID == "" {
		return fmt.Errorf("required field is not set: id")
	}

	if expectedModelId != "" {
		if manifest.ID != expectedModelId {
			return fmt.Errorf("model directory name should match id in manifest: %s != %s", expectedModelId, manifest.ID)
		}
	}

	if manifest.Name == "" {
		return fmt.Errorf("required field is not set: name")
	}

	if manifest.Description == "" {
		return fmt.Errorf("required field is not set: description")
	}

	if manifest.ModelCardUrl == "" {
		return fmt.Errorf("required field is not set: model-card-url")
	}

	if manifest.Quantization == "" {
		return fmt.Errorf("required field is not set: quantization")
	}

	if len(manifest.Capabilities) == 0 {
		return fmt.Errorf("required field is not set: capabilities")
	}

	if manifest.DiskSize == "" {
		return fmt.Errorf("required field is not set: disk-size")
	}

	if len(manifest.Components) == 0 {
		return fmt.Errorf("required field is not set: components")
	}

	if len(manifest.Environment) == 0 {
		return fmt.Errorf("required field is not set: environment")
	}

	return nil
}
