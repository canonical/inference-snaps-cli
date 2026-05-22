package common

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/runtimes"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

// writeRuntimeYAML replaces the runtime.yaml for the named runtime inside the given runtimesDir.
func writeRuntimeYAML(t *testing.T, runtimesDir, name, content string) {
	t.Helper()
	dir := filepath.Join(runtimesDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "runtime.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile runtime.yaml: %v", err)
	}
}

func TestServerEndpoints(t *testing.T) {
	testCases := []struct {
		name            string
		runtimeYAML     string
		wantEndpoints   map[string]string
		wantErrContains string
	}{
		{
			name: "openai http server",
			runtimeYAML: `servers:
  openai:
    protocol: http
    base-path: /v1
`,
			wantEndpoints: map[string]string{
				"openai": "http://localhost:8080/v1",
			},
		},
		{
			name: "unsupported protocol",
			runtimeYAML: `servers:
  openai:
    protocol: ftp
`,
			wantErrContains: "unsupported protocol",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dirs := setupComponentTestDirs(t)
			writeRuntimeYAML(t, dirs.runtimesDir, "test-runtime", tc.runtimeYAML)

			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			cache := storage.NewMockCache()
			_ = cache.SetActiveEngine("test-engine")

			ctx := &Context{
				EnginesDir:  dirs.enginesDir,
				RuntimesDir: dirs.runtimesDir,
				Config:      config,
				Cache:       cache,
			}

			got, err := ServerEndpoints(ctx)
			if tc.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrContains)
				}
				if !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Fatalf("got error %q, want it to contain %q", err.Error(), tc.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			for name, wantURL := range tc.wantEndpoints {
				gotURL, found := got[name]
				if !found {
					t.Fatalf("missing endpoint %q", name)
				}
				if gotURL != wantURL {
					t.Fatalf("endpoint %q: got %q, want %q", name, gotURL, wantURL)
				}
			}
		})
	}
}

func TestServerHttpUrl(t *testing.T) {
	testCases := []struct {
		name   string
		server runtimes.Server
		want   string
	}{
		{
			name:   "default base path",
			server: runtimes.Server{Protocol: "http"},
			want:   "http://localhost:8080/",
		},
		{
			name:   "custom base path",
			server: runtimes.Server{Protocol: "http", BasePath: "/v1"},
			want:   "http://localhost:8080/v1",
		},
		{
			name:   "https protocol",
			server: runtimes.Server{Protocol: "https", BasePath: "/v3"},
			want:   "https://localhost:8080/v3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			ctx := &Context{Config: config}

			got, err := serverHttpUrl(ctx, tc.server)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
