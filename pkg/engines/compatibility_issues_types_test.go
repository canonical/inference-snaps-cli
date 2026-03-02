package engines

import (
	"testing"
)

func TestMarshalYAML(t *testing.T) {

	t.Run("DeviceCompatibilityIssue yaml", func(t *testing.T) {
		DeviceCompatibilityIssue := DeviceCompatibilityIssue{}

		yamlBytes, err := DeviceCompatibilityIssue.MarshalYAML()
		if err != nil {
			t.Fatalf("Failed to marshal DeviceCompatibilityIssue to YAML: %v", err)
		}

		expectedYaml := "required device not found"
		yamlStr, ok := yamlBytes.(string)
		if !ok {
			t.Fatalf("Failed to assert yamlBytes to string")
		}
		if yamlStr != expectedYaml {
			t.Fatalf("Expected YAML: %s, got: %s", expectedYaml, yamlStr)
		}
	})

	t.Run("MemoryCompatibilityIssue yaml", func(t *testing.T) {
		MemoryCompatibilityIssue := MemoryCompatibilityIssue{
			RequiredMemory:  8 * 1024 * 1024 * 1024, // 8 GiB
			AvailableMemory: 4 * 1024 * 1024 * 1024, // 4 GiB
		}

		yamlBytes, err := MemoryCompatibilityIssue.MarshalYAML()
		if err != nil {
			t.Fatalf("Failed to marshal MemoryCompatibilityIssue to YAML: %v", err)
		}

		expectedYaml := "insufficient memory"
		yamlStr, ok := yamlBytes.(string)
		if !ok {
			t.Fatalf("Failed to assert yamlBytes to string")
		}
		if yamlStr != expectedYaml {
			t.Fatalf("Expected YAML: %s, got: %s", expectedYaml, yamlStr)
		}
	})

	t.Run("DiskCompatibilityIssue yaml", func(t *testing.T) {
		DiskCompatibilityIssue := DiskCompatibilityIssue{
			RequiredSpace:  100 * 1024 * 1024 * 1024, // 100 GiB
			AvailableSpace: 50 * 1024 * 1024 * 1024,  // 50 GiB
		}

		yamlBytes, err := DiskCompatibilityIssue.MarshalYAML()
		if err != nil {
			t.Fatalf("Failed to marshal DiskCompatibilityIssue to YAML: %v", err)
		}

		expectedYaml := "insufficient disk space"
		yamlStr, ok := yamlBytes.(string)
		if !ok {
			t.Fatalf("Failed to assert yamlBytes to string")
		}
		if yamlStr != expectedYaml {
			t.Fatalf("Expected YAML: %s, got: %s", expectedYaml, yamlStr)
		}
	})
}

func TestMarshalJson(t *testing.T) {

	t.Run("DeviceCompatibilityIssue json", func(t *testing.T) {
		DeviceCompatibilityIssue := DeviceCompatibilityIssue{}

		jsonBytes, err := DeviceCompatibilityIssue.MarshalJSON()
		if err != nil {
			t.Fatalf("Failed to marshal DeviceCompatibilityIssue to JSON: %v", err)
		}

		expectedJson := "\"required device not found\""
		if string(jsonBytes) != expectedJson {
			t.Fatalf("Expected JSON: %s, got: %s", expectedJson, string(jsonBytes))
		}
	})

	t.Run("MemoryCompatibilityIssue json", func(t *testing.T) {
		MemoryCompatibilityIssue := MemoryCompatibilityIssue{
			RequiredMemory:  8 * 1024 * 1024 * 1024, // 8 GiB
			AvailableMemory: 4 * 1024 * 1024 * 1024, // 4 GiB
		}

		jsonBytes, err := MemoryCompatibilityIssue.MarshalJSON()
		if err != nil {
			t.Fatalf("Failed to marshal MemoryCompatibilityIssue to JSON: %v", err)
		}

		expectedJson := "\"insufficient memory\""
		if string(jsonBytes) != expectedJson {
			t.Fatalf("Expected JSON: %s, got: %s", expectedJson, string(jsonBytes))
		}
	})

	t.Run("DiskCompatibilityIssue json", func(t *testing.T) {
		DiskCompatibilityIssue := DiskCompatibilityIssue{
			RequiredSpace:  100 * 1024 * 1024 * 1024, // 100 GiB
			AvailableSpace: 50 * 1024 * 1024 * 1024,  // 50 GiB
		}

		jsonBytes, err := DiskCompatibilityIssue.MarshalJSON()
		if err != nil {
			t.Fatalf("Failed to marshal DiskCompatibilityIssue to JSON: %v", err)
		}

		expectedJson := "\"insufficient disk space\""
		if string(jsonBytes) != expectedJson {
			t.Fatalf("Expected JSON: %s, got: %s", expectedJson, string(jsonBytes))
		}
	})
}
