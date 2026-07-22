package common

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func TestServerListeners(t *testing.T) {
	testCases := []struct {
		name             string
		engineYAML       string
		runtimeYAML      string
		wantListenerURLs map[string]string
		wantErrContains  string
	}{
		{
			name: "multiple servers",
			engineYAML: `name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `servers:
  openai:
    protocol: http
    base-path: /v1
  kserve:
    protocol: https
    base-path: /v2
`,
			wantListenerURLs: map[string]string{
				"openai": "http://127.0.0.1:8080/v1",
				"kserve": "https://127.0.0.1:8080/v2",
				"webui":  "http://192.0.2.1:8080/",
			},
		},
		{
			name: "unsupported protocol",
			engineYAML: `name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `servers:
  openai:
    protocol: ftp
`,
			wantErrContains: "unsupported protocol",
		},
		{
			name: "websocket servers",
			engineYAML: `name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `servers:
  openai-ws:
    protocol: ws
    base-path: /v1
  openai-unix:
    protocol: ws+unix
    base-path: /v1
`,
			wantListenerURLs: map[string]string{
				"openai-ws":   "ws://127.0.0.1:8081/v1",
				"openai-unix": "ws://127.0.0.1:8081/v1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ADDITIONAL_FEATURES", "webui")

			enginesDir := t.TempDir()
			runtimesDir := t.TempDir()

			const (
				engineName  = "test-engine"
				runtimeName = "test-runtime"
			)

			writeEngineYAML(t, enginesDir, engineName, tc.engineYAML)
			writeRuntimeYAML(t, runtimesDir, runtimeName, tc.runtimeYAML)

			cache := storage.NewMockCache()
			if err := cache.SetActiveEngine(engineName); err != nil {
				t.Fatalf("SetActiveEngine: %v", err)
			}

			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			config.Set("http.host", "127.0.0.1", storage.UserConfig)
			config.Set("webui.http.port", "8080", storage.UserConfig)
			config.Set("webui.http.host", "192.0.2.1", storage.UserConfig)
			config.Set("ws.port", "8081", storage.UserConfig)
			config.Set("ws.host", "127.0.0.1", storage.UserConfig)
			config.Set("ws.unix-socket", "/run/openai.sock", storage.UserConfig)

			ctx := &Context{
				EnginesDir:  enginesDir,
				RuntimesDir: runtimesDir,
				Config:      config,
				Cache:       cache,
			}

			got, err := ServerListeners(ctx)
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
			for name, wantURL := range tc.wantListenerURLs {
				listener, found := got[name]
				if !found {
					t.Fatalf("missing listener %q", name)
				}
				if listener.Url != wantURL {
					t.Fatalf("listener %q: got %q, want %q", name, listener.Url, wantURL)
				}
			}
		})
	}
}

func TestServerHTTPListener(t *testing.T) {
	testCases := []struct {
		name    string
		server  runtimes.Server
		host    string
		setHost bool
		want    string
		wantErr bool
	}{
		{
			name: "default base path",
			server: runtimes.Server{
				Protocol: "http",
			},
			host:    "0.0.0.0",
			setHost: true,
			want:    "http://0.0.0.0:8080/",
		},
		{
			name: "custom base path",
			server: runtimes.Server{
				Protocol: "http",
				BasePath: "/v1",
			},
			host:    "127.0.0.1",
			setHost: true,
			want:    "http://127.0.0.1:8080/v1",
		},
		{
			name: "https protocol",
			server: runtimes.Server{
				Protocol: "https",
				BasePath: "/v3",
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
			if tc.setHost {
				config.Set("http.host", tc.host, storage.UserConfig)
			}
			ctx := &Context{
				Config: config,
			}

			got, err := serverHttpListener(ctx, tc.server)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("expected non-nil listener")
			}
			if got.Url != tc.want {
				t.Fatalf("got %q, want %q", got.Url, tc.want)
			}
		})
	}
}

func TestServerWSListener(t *testing.T) {
	testCases := []struct {
		name    string
		server  runtimes.Server
		host    string
		setHost bool
		want    string
	}{
		{
			name: "default base path",
			server: runtimes.Server{
				Protocol: "ws",
			},
			host:    "0.0.0.0",
			setHost: true,
			want:    "ws://0.0.0.0:8080/",
		},
		{
			name: "custom base path",
			server: runtimes.Server{
				Protocol: "ws",
				BasePath: "/v1",
			},
			host:    "127.0.0.1",
			setHost: true,
			want:    "ws://127.0.0.1:8080/v1",
		},
		{
			name: "complex base path",
			server: runtimes.Server{
				Protocol: "ws",
				BasePath: "/api/v2",
			},
			host:    "localhost",
			setHost: true,
			want:    "ws://localhost:8080/api/v2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			config.Set("ws.port", "8080", storage.UserConfig)
			if tc.setHost {
				config.Set("ws.host", tc.host, storage.UserConfig)
			}
			ctx := &Context{
				Config: config,
			}

			got, err := serverWsListener(ctx, "test-server", tc.server)
			if err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("expected non-nil listener")
			}
			if got.Url != tc.want {
				t.Fatalf("got %q, want %q", got.Url, tc.want)
			}
		})
	}
}

func TestServerWsUnixListener(t *testing.T) {
	testCases := []struct {
		name       string
		server     runtimes.Server
		host       string
		setHost    bool
		wantURL    string
		wantSocket string
	}{
		{
			name: "default base path",
			server: runtimes.Server{
				Protocol: "ws+unix",
			},
			host:       "0.0.0.0",
			setHost:    true,
			wantURL:    "ws://0.0.0.0:8080/",
			wantSocket: "/run/socket.sock",
		},
		{
			name: "custom base path",
			server: runtimes.Server{
				Protocol: "ws+unix",
				BasePath: "/v1",
			},
			host:       "127.0.0.1",
			setHost:    true,
			wantURL:    "ws://127.0.0.1:8080/v1",
			wantSocket: "/run/socket.sock",
		},
		{
			name: "complex base path",
			server: runtimes.Server{
				Protocol: "ws+unix",
				BasePath: "/api/stream",
			},
			host:       "localhost",
			setHost:    true,
			wantURL:    "ws://localhost:8080/api/stream",
			wantSocket: "/run/snap.whisper/openai.sock",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			config.Set("ws.port", "8080", storage.UserConfig)
			if tc.setHost {
				config.Set("ws.host", tc.host, storage.UserConfig)
			}
			config.Set("ws.unix-socket", tc.wantSocket, storage.UserConfig)

			ctx := &Context{
				Config: config,
			}

			got, err := serverWsOverUnixListener(ctx, "test-server", tc.server)
			if err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("expected non-nil listener")
			}
			if got.Url != tc.wantURL {
				t.Fatalf("URL: got %q, want %q", got.Url, tc.wantURL)
			}
			if got.UnixSocket != tc.wantSocket {
				t.Fatalf("UnixSocket: got %q, want %q", got.UnixSocket, tc.wantSocket)
			}
		})
	}
}
