package common

import (
	"testing"
)

func TestModelPropertiesEmpty(t *testing.T) {
	result, err := modelProperties(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "" || result.Quantization != "" || result.MultimediaProjector != "" {
		t.Errorf("expected empty struct, got %+v", result)
	}
}

func TestModelPropertiesAllFields(t *testing.T) {
	settings := []ComponentSettings{
		{
			Properties: map[string]string{
				"model-name":         "llama3",
				"model-quantization": "q4_k_m",
				"mmproj-name":        "mmproj-clip",
			},
		},
	}

	result, err := modelProperties(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "llama3" {
		t.Errorf("expected Name %q, got %q", "llama3", result.Name)
	}
	if result.Quantization != "q4_k_m" {
		t.Errorf("expected Quantization %q, got %q", "q4_k_m", result.Quantization)
	}
	if result.MultimediaProjector != "mmproj-clip" {
		t.Errorf("expected MultimediaProjector %q, got %q", "mmproj-clip", result.MultimediaProjector)
	}
}

func TestModelPropertiesPartialFields(t *testing.T) {
	settings := []ComponentSettings{
		{
			Properties: map[string]string{
				"model-name": "mymodel",
			},
		},
	}

	result, err := modelProperties(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "mymodel" {
		t.Errorf("expected Name %q, got %q", "mymodel", result.Name)
	}
	if result.Quantization != "" {
		t.Errorf("expected empty Quantization, got %q", result.Quantization)
	}
	if result.MultimediaProjector != "" {
		t.Errorf("expected empty MultimediaProjector, got %q", result.MultimediaProjector)
	}
}

func TestModelPropertiesNoRelevantProperties(t *testing.T) {
	settings := []ComponentSettings{
		{
			Properties: map[string]string{
				"some-other-key": "some-value",
			},
		},
	}

	result, err := modelProperties(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "" || result.Quantization != "" || result.MultimediaProjector != "" {
		t.Errorf("expected empty struct, got %+v", result)
	}
}

func TestModelPropertiesMultipleComponentsLaterOverrides(t *testing.T) {
	settings := []ComponentSettings{
		{
			Properties: map[string]string{
				"model-name":         "first-model",
				"model-quantization": "q4",
			},
		},
		{
			Properties: map[string]string{
				"model-name":  "second-model",
				"mmproj-name": "mmproj-v2",
			},
		},
	}

	result, err := modelProperties(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "second-model" {
		t.Errorf("expected Name %q, got %q", "second-model", result.Name)
	}
	if result.Quantization != "q4" {
		t.Errorf("expected Quantization %q, got %q", "q4", result.Quantization)
	}
	if result.MultimediaProjector != "mmproj-v2" {
		t.Errorf("expected MultimediaProjector %q, got %q", "mmproj-v2", result.MultimediaProjector)
	}
}

func TestModelPropertiesNilPropertiesMap(t *testing.T) {
	settings := []ComponentSettings{
		{
			Properties: nil,
		},
	}

	result, err := modelProperties(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "" || result.Quantization != "" || result.MultimediaProjector != "" {
		t.Errorf("expected empty struct, got %+v", result)
	}
}
