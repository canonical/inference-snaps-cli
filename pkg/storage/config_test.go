package storage

import (
	"os"
	"testing"
)

func TestSetEngineOverridesPackage(t *testing.T) {
	config := NewMockConfig()

	// Set environment config
	os.Setenv("DEFAULT_API_ENDPOINT", "https://default.example.com")

	// Set engine config with different value
	err := config.Set("model", "mistral", EngineConfig)
	if err != nil {
		t.Fatalf("Set engine config returned unexpected error: %v", err)
	}

	// Verify engine value overrides package value
	values, err := config.Get("model")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}

	if value, found := values["model"]; !found || value != "mistral" {
		t.Fatalf("expected engine config to override package config, got %#v", values)
	}
}

func TestSetUserOverridesEngine(t *testing.T) {
	config := NewMockConfig()

	// Set engine config
	err := config.Set("model", "mistral", EngineConfig)
	if err != nil {
		t.Fatalf("Set engine config returned unexpected error: %v", err)
	}

	// Set user config with different value
	err = config.Set("model", "custom-model", UserConfig)
	if err != nil {
		t.Fatalf("Set user config returned unexpected error: %v", err)
	}

	// Verify user value overrides engine value
	values, err := config.Get("model")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}

	if value, found := values["model"]; !found || value != "custom-model" {
		t.Fatalf("expected user config to override engine config, got %#v", values)
	}
}

func TestSetUserOverridesPackage(t *testing.T) {
	config := NewMockConfig()

	// Set environment config
	os.Setenv("DEFAULT_API_ENDPOINT", "https://default.example.com")

	// Set user config with different value (skipping engine level)
	err := config.Set("api-endpoint", "https://user.example.com", UserConfig)
	if err != nil {
		t.Fatalf("Set user config returned unexpected error: %v", err)
	}

	// Verify user value overrides package value
	values, err := config.Get("api-endpoint")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}

	if value, found := values["api-endpoint"]; !found || value != "https://user.example.com" {
		t.Fatalf("expected user config to override package config, got %#v", values)
	}
}

func TestGetAll(t *testing.T) {
	config := NewMockConfig()
	config.Set("model-override", "engine", EngineConfig)
	config.Set("model-override", "custom-model", UserConfig)
	os.Setenv("DEFAULT_HTTP_PORT", "8080")
	os.Setenv("DEFAULT_HTTP_HOST", "8080")
	config.Set("http-port", "8000", EngineConfig)
	config.Set("http-port", "8001", UserConfig)
	v, err := config.GetAll()
	if err != nil {
		t.Fatalf("GetAll returned unexpected error: %v", err)
	}
	if v["model-override"] != "custom-model" {
		t.Fatalf("expected model-override to be custom-model, got %#v", v["model-override"])
	}
	if v["http-port"] != "8001" {
		t.Fatalf("expected http-port to be overridden by user config to 8001, got %#v", v["http-port"])
	}
	if v["http-host"] != "8080" {
		t.Fatalf("expected http-host to be overridden by environment variable to 8080, got %#v", v["http-host"])
	}
}
