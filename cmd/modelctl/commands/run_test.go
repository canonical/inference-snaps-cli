package commands

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func testRunContext(t *testing.T, runtimeYAML string) *common.Context {
	t.Helper()

	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	runtimesDir := filepath.Join(base, "runtimes")
	if err := os.MkdirAll(filepath.Join(enginesDir, "test-engine"), 0o755); err != nil {
		t.Fatalf("creating engine dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(runtimesDir, "test-runtime"), 0o755); err != nil {
		t.Fatalf("creating runtime dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(enginesDir, "test-engine", "engine.yaml"), []byte("name: test-engine\nruntime: test-runtime\n"), 0o644); err != nil {
		t.Fatalf("writing engine manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runtimesDir, "test-runtime", "runtime.yaml"), []byte(runtimeYAML), 0o644); err != nil {
		t.Fatalf("writing runtime manifest: %v", err)
	}

	cfg := storage.NewMockConfig()
	if err := cfg.Set("http.host", "127.0.0.1", storage.UserConfig); err != nil {
		t.Fatalf("setting http.host: %v", err)
	}
	if err := cfg.Set("http.port", "8080", storage.UserConfig); err != nil {
		t.Fatalf("setting http.port: %v", err)
	}

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("test-engine"); err != nil {
		t.Fatalf("setting active engine: %v", err)
	}

	return &common.Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		Cache:       cache,
		Config:      cfg,
	}
}

func TestProcessEnvConfigs(t *testing.T) {
	mockConfig := storage.NewMockConfig()
	mockConfig.Set("env.my-key", "value", storage.UserConfig)
	mockConfig.Set("env.other", "123", storage.UserConfig)
	mockConfig.Set("other.ignored", "ignored", storage.UserConfig)
	cmd := runCommand{
		Context: &common.Context{
			Config: mockConfig,
		},
	}

	err := cmd.processEnvConfigs()
	if err != nil {
		t.Fatalf("processEnvConfigs returned error: %v", err)
	}

	if got := os.Getenv("MY_KEY"); got != "value" {
		t.Fatalf("expected MY_KEY to be %q, got %q", "value", got)
	}

	if got := os.Getenv("OTHER"); got != "123" {
		t.Fatalf("expected OTHER to be %q, got %q", "123", got)
	}
}

func TestWriteShareProviderEnv(t *testing.T) {
	t.Run("default path", func(t *testing.T) {
		t.Setenv("SNAP_COMMON", t.TempDir())
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")
		expectedProviderDir := os.ExpandEnv("$SNAP_COMMON/share/provider")

		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")}
		cmd.shareProvider = cmd.defaultProviderDirectoryPath()
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		cachedShareProviderDirectory, err := cmd.Cache.GetSharedProviderDirectory()
		if err != nil {
			t.Fatalf("getting cached shared provider directory: %v", err)
		}
		if cachedShareProviderDirectory != expectedProviderDir {
			t.Fatalf("cached shared provider directory mismatch\nwant: %q\ngot:  %q", expectedProviderDir, cachedShareProviderDirectory)
		}

		path := filepath.Join(expectedProviderDir, "provider.env")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nOPENAI_BASE_URL=http://127.0.0.1:8080/v1\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	t.Run("custom path", func(t *testing.T) {
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		path := filepath.Join(t.TempDir(), "custom")
		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n"), shareProvider: path}
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		cachedShareProviderDirectory, err := cmd.Cache.GetSharedProviderDirectory()
		if err != nil {
			t.Fatalf("getting cached shared provider directory: %v", err)
		}
		if cachedShareProviderDirectory != path {
			t.Fatalf("cached shared provider directory mismatch\nwant: %q\ngot:  %q", path, cachedShareProviderDirectory)
		}

		content, err := os.ReadFile(filepath.Join(path, "provider.env"))
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nOPENAI_BASE_URL=http://127.0.0.1:8080/v1\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	t.Run("runtime without openai entry", func(t *testing.T) {
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		path := t.TempDir()
		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  kserve:\n    protocol: http\n    base-path: /v2\n"), shareProvider: path}
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		cachedShareProviderDirectory, err := cmd.Cache.GetSharedProviderDirectory()
		if err != nil {
			t.Fatalf("getting cached shared provider directory: %v", err)
		}
		if cachedShareProviderDirectory != path {
			t.Fatalf("cached shared provider directory mismatch\nwant: %q\ngot:  %q", path, cachedShareProviderDirectory)
		}

		content, err := os.ReadFile(filepath.Join(path, "provider.env"))
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

}

func TestNoActiveModel(t *testing.T) {
	t.Run("no active model", func(t *testing.T) {
		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")}
		err := cmd.run(nil, []string{"echo", "Hello World!"})
		if err == nil || err.Error() != "no active model" {
			t.Fatalf("expected error 'no active model', got %v", err)
		}
	})
}

func TestRunCommandFailure(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "non-zero exit", args: []string{"/bin/sh", "-c", "exit 7"}},
		{name: "executable not found", args: []string{filepath.Join(t.TempDir(), "missing-command")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")
			modelsDir := t.TempDir()
			modelDir := filepath.Join(modelsDir, "test-model")
			if err := os.Mkdir(modelDir, 0o755); err != nil {
				t.Fatalf("creating model dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(modelDir, "model.yaml"), []byte("name: test-model\n"), 0o644); err != nil {
				t.Fatalf("writing model manifest: %v", err)
			}
			ctx.ModelsDir = modelsDir
			if err := ctx.Cache.SetActiveModel("test-model"); err != nil {
				t.Fatalf("setting active model: %v", err)
			}

			err := (&runCommand{Context: ctx}).run(nil, tt.args)
			if err == nil || !strings.HasPrefix(err.Error(), "command exited with non-zero status: ") {
				t.Fatalf("expected friendly command failure, got %v", err)
			}
		})
	}
}

func TestCommandStopped(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{name: "SIGTERM", command: "kill -TERM $$", want: true},
		{name: "SIGINT", command: "kill -INT $$", want: true},
		{name: "SIGKILL", command: "kill -KILL $$", want: false},
		{name: "non-zero exit", command: "exit 7", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exec.Command("/bin/sh", "-c", tt.command).Run()
			if err == nil {
				t.Fatal("expected command to fail")
			}
			if got := commandStopped(err); got != tt.want {
				var status syscall.WaitStatus
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					status, _ = exitErr.ProcessState.Sys().(syscall.WaitStatus)
				}
				t.Fatalf("commandStopped() = %v, want %v (status %v)", got, tt.want, status)
			}
		})
	}
}
