package common

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestServerHTTPEndpoints(t *testing.T) {
	testCases := []struct {
		name              string
		httpEndpointNames string
		env               map[string]string
		want              map[string]string
	}{
		{
			name:              "uses default base path when endpoint base path env var is absent",
			httpEndpointNames: "openai",
			env:               map[string]string{},
			want: map[string]string{
				"openai": "http://localhost:8080/",
			},
		},
		{
			name:              "uses endpoint specific base path from env var",
			httpEndpointNames: "openai",
			env: map[string]string{
				"HTTP_OPENAI_BASE_PATH": "/v1",
			},
			want: map[string]string{
				"openai": "http://localhost:8080/v1",
			},
		},
		{
			name:              "trims endpoint names and supports multiple endpoints",
			httpEndpointNames: " openai , webui ",
			env: map[string]string{
				"HTTP_OPENAI_BASE_PATH": "/v1",
			},
			want: map[string]string{
				"openai": "http://localhost:8080/v1",
				"webui":  "http://localhost:8080/",
			},
		},
		// TODO test component without any endpoint
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			ctx := &Context{
				Config: storage.NewMockConfig(map[string]any{"http.port": "8080"}),
			}

			got, err := serverHTTPEndpoints(ctx, testCase.httpEndpointNames)
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
