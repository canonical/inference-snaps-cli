package common

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestServerEndpoints(t *testing.T) {
	testCases := []struct {
		name             string
		componentConfigs []ComponentSettings
		want             map[string]string
		wantErrContains  string
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
				"openai": "http://localhost:8080/v1",
				"kserve": "https://localhost:8080/v2",
				"webui":  "http://localhost:8080/",
			},
		},
		{
			name: "unsupported protocol",
			componentConfigs: []ComponentSettings{
				{
					Servers: map[string]map[string]string{
						"openai": {
							"protocol": "ftp",
						},
					},
				},
			},
			wantErrContains: "unsupported protocol",
		},
		{
			name: "OPENAI_BASE_PATH deprecated",
			componentConfigs: []ComponentSettings{
				{
					Environment: []string{"OPENAI_BASE_PATH=/v1"},
				},
			},
			wantErrContains: "deprecated",
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
		configValues map[string]any
		serverConfig map[string]string
		want         string
	}{
		{
			name: "default base path",
			configValues: map[string]any{
				"http.port": "8080",
			},
			serverConfig: map[string]string{
				"protocol": "http",
			},
			want: "http://localhost:8080/",
		},
		{
			name: "custom base path",
			configValues: map[string]any{
				"http.port": "8080",
			},
			serverConfig: map[string]string{
				"protocol":  "http",
				"base-path": "/v1",
			},
			want: "http://localhost:8080/v1",
		},
		{
			name: "https protocol",
			configValues: map[string]any{
				"http.port": "8080",
			},
			serverConfig: map[string]string{
				"protocol":  "https",
				"base-path": "/v3",
			},
			want: "https://localhost:8080/v3",
		},
		{
			name: "configured fqdn",
			configValues: map[string]any{
				"http.port": "8080",
				"http.fqdn": "inference.example.com",
			},
			serverConfig: map[string]string{
				"protocol":  "http",
				"base-path": "/v1",
			},
			want: "http://inference.example.com:8080/v1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := storage.NewMockConfig(testCase.configValues)
			ctx := &Context{
				Config: config,
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
