package commands

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

// helper to create a test context with test data
func createTestContextForStatus(t *testing.T) *common.Context {
	t.Helper()
	testDataDir := "../../../test_data"
	enginesDir := filepath.Join(testDataDir, "engines")
	runtimesDir := filepath.Join(testDataDir, "runtimes")
	modelsDir := filepath.Join(testDataDir, "models")

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("cpu"); err != nil {
		t.Fatalf("failed to set active engine: %v", err)
	}
	if err := cache.SetActiveModel("4b-it-int4-fq-ov"); err != nil {
		t.Fatalf("failed to set active model: %v", err)
	}

	config := storage.NewMockConfig()
	if err := config.Set("http.port", "8080", storage.UserConfig); err != nil {
		t.Fatalf("failed to set http.port: %v", err)
	}
	if err := config.Set("http.host", "127.0.0.1", storage.UserConfig); err != nil {
		t.Fatalf("failed to set http.host: %v", err)
	}
	if err := config.Set("ws.unix-socket", "/run/whisper.sock", storage.UserConfig); err != nil {
		t.Fatalf("failed to set ws.unix-socket: %v", err)
	}

	return &common.Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   modelsDir,
		Cache:       cache,
		Config:      config,
		Snap:        snap.MockWithServiceStatuses(map[string]string{"llama-server": "active"}),
	}
}

// createTestContextForStatusExample creates a test context for Example tests
func createTestContextForStatusExample() *common.Context {
	testDataDir := "../../../test_data"
	enginesDir := filepath.Join(testDataDir, "engines")
	runtimesDir := filepath.Join(testDataDir, "runtimes")
	modelsDir := filepath.Join(testDataDir, "models")

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("cpu"); err != nil {
		panic(fmt.Sprintf("failed to set active engine: %v", err))
	}
	if err := cache.SetActiveModel("4b-it-int4-fq-ov"); err != nil {
		panic(fmt.Sprintf("failed to set active model: %v", err))
	}

	config := storage.NewMockConfig()
	if err := config.Set("http.port", "8080", storage.UserConfig); err != nil {
		panic(fmt.Sprintf("failed to set http.port: %v", err))
	}
	if err := config.Set("http.host", "127.0.0.1", storage.UserConfig); err != nil {
		panic(fmt.Sprintf("failed to set http.host: %v", err))
	}
	if err := config.Set("ws.unix-socket", "/run/whisper.sock", storage.UserConfig); err != nil {
		panic(fmt.Sprintf("failed to set ws.unix-socket: %v", err))
	}

	return &common.Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   modelsDir,
		Cache:       cache,
		Config:      config,
		Snap:        snap.MockWithServiceStatuses(map[string]string{"llama-server": "active"}),
	}
}

func Example_statusCommand_printStatusYaml() {
	ctx := createTestContextForStatusExample()
	cmd := statusCommand{
		Context: ctx,
		format:  "yaml",
	}

	if err := cmd.run(nil, nil); err != nil {
		panic(fmt.Sprintf("failed to print status in yaml format: %v", err))
	}

	// Output:
	// engine: cpu
	// services:
	//     llama-server: active
	// entrypoints:
	//     openai:
	//         url: http://127.0.0.1:8080/v1
	//     whisperlive:
	//         unix-socket: '/run/whisper.sock (url: ws://unix/realtime)'
	// model:
	//     name: gemma-4-4b-it-int4-fq-ov
}

func Example_statusCommand_printStatusJson() {
	ctx := createTestContextForStatusExample()
	cmd := statusCommand{
		Context: ctx,
		format:  "json",
	}

	if err := cmd.run(nil, nil); err != nil {
		panic(fmt.Sprintf("failed to print status in json format: %v", err))
	}

	// Output:
	// {
	//   "engine": "cpu",
	//   "services": {
	//     "llama-server": "active"
	//   },
	//   "entrypoints": {
	//     "openai": {
	//       "url": "http://127.0.0.1:8080/v1"
	//     },
	//     "whisperlive": {
	//       "unix-socket": "/run/whisper.sock",
	//       "unix-socket-url": "ws://unix/realtime"
	//     }
	//   },
	//   "model": {
	//     "name": "gemma-4-4b-it-int4-fq-ov"
	//   }
	// }
}

func TestStatusCommand_statusYaml(t *testing.T) {
	ctx := createTestContextForStatus(t)
	cmd := &statusCommand{Context: ctx}

	yamlStr, err := cmd.statusYaml()
	if err != nil {
		t.Fatalf("statusYaml() returned error: %v", err)
	}

	if yamlStr == "" {
		t.Error("expected non-empty YAML output")
	}
}

func TestStatusCommand_statusJson(t *testing.T) {
	ctx := createTestContextForStatus(t)
	cmd := &statusCommand{Context: ctx}

	jsonStr, err := cmd.statusJson()
	if err != nil {
		t.Fatalf("statusJson() returned error: %v", err)
	}

	if jsonStr == "" {
		t.Error("expected non-empty JSON output")
	}
}
