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
			name: "multiple components and servers",
			componentConfigs: []ComponentSettings{
				{
					Servers: map[string]map[string]string{
						"openai": {
							"protocol":  "http",
							"base-path": "/v1",
						},
					},
				},
				{
					Servers: map[string]map[string]string{
						"kserve": {
							"protocol":  "https",
							"base-path": "/v2",
						},
						"webui": {
							"protocol": "http",
						},
					},
				},
			},
			want: map[string]string{
				"openai": "http://127.0.0.1:8080/v1",
				"kserve": "https://127.0.0.1:8080/v2",
				"webui":  "http://192.0.2.1:8080/",
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

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("ADDITIONAL_FEATURES", "webui")

			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			config.Set("http.host", "127.0.0.1", storage.UserConfig)
			config.Set("webui.http.port", "8080", storage.UserConfig)
			config.Set("webui.http.host", "192.0.2.1", storage.UserConfig)
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
		name         string
		serverConfig map[string]string
		host         string
		setHost      bool
		want         string
	}{
		{
			name: "default base path",
			serverConfig: map[string]string{
				"protocol": "http",
			},
			host:    "0.0.0.0",
			setHost: true,
			want:    "http://0.0.0.0:8080/",
		},
		{
			name: "custom base path",
			serverConfig: map[string]string{
				"protocol":  "http",
				"base-path": "/v1",
			},
			host:    "127.0.0.1",
			setHost: true,
			want:    "http://127.0.0.1:8080/v1",
		},
		{
			name: "https protocol",
			serverConfig: map[string]string{
				"protocol":  "https",
				"base-path": "/v3",
			},
			host:    "0.0.0.0",
			setHost: true,
			want:    "https://0.0.0.0:8080/v3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			if testCase.setHost {
				config.Set("http.host", testCase.host, storage.UserConfig)
			}
			ctx := &Context{
				Config: config,
			}

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
