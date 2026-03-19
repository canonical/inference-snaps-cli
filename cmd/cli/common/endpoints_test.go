package common

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestServerEndpoints(t *testing.T) {
	testCases := []struct {
		name             string
		componentConfigs []ComponentConfig
		want             map[string]string
		wantErrContains  string
	}{
		{
			name: "builds endpoints from multiple components",
			componentConfigs: []ComponentConfig{
				{
					Servers: map[string]map[string]string{
						"openai": {
							"http":      "http",
							"base-path": "/v1",
						},
					},
				},
				{
					Servers: map[string]map[string]string{
						"kserve": {
							"http":      "https",
							"base-path": "/v2",
						},
						"webui": {
							"http": "http",
						},
					},
				},
			},
			want: map[string]string{
				"openai": "http://localhost:8080/v1",
				"kserve": "https://localhost:8080/v2",
				"webui":  "http://localhost:8080/",
			},
		},
		{
			name: "returns error for unsupported protocol",
			componentConfigs: []ComponentConfig{
				{
					Servers: map[string]map[string]string{
						"openai": {
							"http":     "ftp",
							"protocol": "ftp",
						},
					},
				},
			},
			wantErrContains: `unsupported protocol "ftp" for server "openai"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := &Context{
				Config: storage.NewMockConfig(map[string]any{"http.port": "8080"}),
			}

			got, err := serverEndpoints(ctx, testCase.componentConfigs)
			if testCase.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", testCase.wantErrContains)
				}
				if !strings.Contains(err.Error(), testCase.wantErrContains) {
					t.Fatalf("got error %q, want it to contain %q", err.Error(), testCase.wantErrContains)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if len(got) != len(testCase.want) {
				t.Fatalf("got %d endpoints, want %d", len(got), len(testCase.want))
			}

			for endpointName, wantURL := range testCase.want {
				gotURL, found := got[endpointName]
				if !found {
					t.Fatalf("missing endpoint %q", endpointName)
				}
				if gotURL != wantURL {
					t.Fatalf("got %q: %q, want %q", endpointName, gotURL, wantURL)
				}
			}
		})
	}
}

func TestServerHttpUrl(t *testing.T) {
	testCases := []struct {
		name         string
		serverConfig map[string]string
		want         string
	}{
		{
			name: "uses default base path when not provided",
			serverConfig: map[string]string{
				"http": "http",
			},
			want: "http://localhost:8080/",
		},
		{
			name: "uses configured base path",
			serverConfig: map[string]string{
				"http":      "http",
				"base-path": "/v1",
			},
			want: "http://localhost:8080/v1",
		},
		{
			name: "supports https protocol",
			serverConfig: map[string]string{
				"http":      "https",
				"base-path": "/v3",
			},
			want: "https://localhost:8080/v3",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := &Context{
				Config: storage.NewMockConfig(map[string]any{"http.port": "8080"}),
			}

			got, err := serverHttpUrl(ctx, testCase.serverConfig)
			if err != nil {
				t.Fatal(err)
			}

			if got != testCase.want {
				t.Fatalf("got %q, want %q", got, testCase.want)
			}
		})
	}
}
