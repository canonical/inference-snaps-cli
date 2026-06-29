package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteShareFiles_writesStatusJson(t *testing.T) {
	dir := t.TempDir()

	status := &sharedStatus{
		Engine:   "cpu",
		Services: map[string]string{"inference": "active"},
	}

	if err := writeShareFiles(status, dir); err != nil {
		t.Fatalf("writeShareFiles returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "status.json"))
	if err != nil {
		t.Fatalf("reading status.json: %v", err)
	}

	var got sharedStatus
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshalling status.json: %v", err)
	}

	if got.Engine != "cpu" {
		t.Errorf("engine = %q, want %q", got.Engine, "cpu")
	}
	if got.Services["inference"] != "active" {
		t.Errorf("services[inference] = %q, want %q", got.Services["inference"], "active")
	}
}

func TestWriteShareFiles_writesOpenaiJsonWhenEndpointPresent(t *testing.T) {
	dir := t.TempDir()

	status := &sharedStatus{
		Engine:    "cpu",
		Services:  map[string]string{},
		Endpoints: map[string]string{"openai": "http://localhost:8080/v1"},
	}

	if err := writeShareFiles(status, dir); err != nil {
		t.Fatalf("writeShareFiles returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "openai.json"))
	if err != nil {
		t.Fatalf("reading openai.json: %v", err)
	}

	var got sharedOpenai
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshalling openai.json: %v", err)
	}

	if got.BaseUrl != "http://localhost:8080/v1" {
		t.Errorf("base_url = %q, want %q", got.BaseUrl, "http://localhost:8080/v1")
	}
}

func TestWriteShareFiles_noOpenaiJsonWhenEndpointAbsent(t *testing.T) {
	dir := t.TempDir()

	status := &sharedStatus{
		Engine:   "cpu",
		Services: map[string]string{},
	}

	if err := writeShareFiles(status, dir); err != nil {
		t.Fatalf("writeShareFiles returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "openai.json")); !os.IsNotExist(err) {
		t.Errorf("expected openai.json to be absent, but it exists (or unexpected error: %v)", err)
	}
}

func TestWriteShareFiles_removesStaleOpenaiJsonWhenEndpointDropped(t *testing.T) {
	dir := t.TempDir()

	// First call: endpoint present → openai.json is written.
	withEndpoint := &sharedStatus{
		Engine:    "cpu",
		Services:  map[string]string{},
		Endpoints: map[string]string{"openai": "http://localhost:8080/v1"},
	}
	if err := writeShareFiles(withEndpoint, dir); err != nil {
		t.Fatalf("first writeShareFiles returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "openai.json")); err != nil {
		t.Fatalf("expected openai.json to exist after first call, got: %v", err)
	}

	// Second call: endpoint gone → openai.json should be removed.
	withoutEndpoint := &sharedStatus{
		Engine:   "cpu",
		Services: map[string]string{},
	}
	if err := writeShareFiles(withoutEndpoint, dir); err != nil {
		t.Fatalf("second writeShareFiles returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "openai.json")); !os.IsNotExist(err) {
		t.Errorf("expected openai.json to be removed after second call, got: %v", err)
	}
}

