package fallbackserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListenAddress(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
		wantErr bool
	}{
		{name: "host port and path", baseURL: "http://127.0.0.1:8080/v1", want: "127.0.0.1:8080"},
		{name: "default port", baseURL: "http://localhost/v3", want: "localhost:80"},
		{name: "IPv6", baseURL: "http://[::1]:9000/api/v1", want: "[::1]:9000"},
		{name: "no path", baseURL: "http://localhost:8080", want: "localhost:8080"},
		{name: "invalid URL", baseURL: "://localhost:8080/v1", wantErr: true},
		{name: "missing host", baseURL: "http:///v1", wantErr: true},
		{name: "HTTPS unsupported", baseURL: "https://localhost:8080/v1", wantErr: true},
		{name: "user information", baseURL: "http://user@localhost:8080/v1", wantErr: true},
		{name: "query", baseURL: "http://localhost:8080/v1?key=value", wantErr: true},
		{name: "fragment", baseURL: "http://localhost:8080/v1#fragment", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listenAddress(tt.baseURL)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("listenAddress(%q) unexpectedly succeeded with %q", tt.baseURL, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("listenAddress(%q) returned error: %v", tt.baseURL, err)
			}
			if got != tt.want {
				t.Fatalf("listenAddress(%q) = %q, want %q", tt.baseURL, got, tt.want)
			}
		})
	}
}

func TestHandlerReportsServiceUnavailable(t *testing.T) {
	type apiError struct {
		Error struct {
			Message string  `json:"message"`
			Type    string  `json:"type"`
			Param   *string `json:"param"`
			Code    string  `json:"code"`
		} `json:"error"`
	}

	errorMessages := []string{
		"Not enough disk space.",
		`Model "example" could not be installed.`,
	}
	responseHandler, err := handler(errorMessages)
	if err != nil {
		t.Fatalf("creating handler: %v", err)
	}

	for _, request := range []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/"},
		{method: http.MethodGet, path: "/v1/models"},
		{method: http.MethodPost, path: "/v3/chat/completions"},
		{method: http.MethodGet, path: "/unknown"},
	} {
		t.Run(request.method+" "+request.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			responseHandler.ServeHTTP(recorder, httptest.NewRequest(request.method, request.path, nil))

			if recorder.Code != http.StatusServiceUnavailable {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want %q", got, "application/json")
			}
			if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("Cache-Control = %q, want %q", got, "no-store")
			}

			var response apiError
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("decoding response: %v", err)
			}
			if want := strings.Join(errorMessages, "\n"); response.Error.Message != want {
				t.Fatalf("error message = %q, want %q", response.Error.Message, want)
			}
			if response.Error.Type != "unavailable_error" {
				t.Fatalf("error type = %q, want %q", response.Error.Type, "unavailable_error")
			}
			if response.Error.Param != nil {
				t.Fatalf("error param = %q, want nil", *response.Error.Param)
			}
			if response.Error.Code != "service_unavailable" {
				t.Fatalf("error code = %q, want %q", response.Error.Code, "service_unavailable")
			}
		})
	}
}

func TestHandlerWithoutErrorMessages(t *testing.T) {
	responseHandler, err := handler(nil)
	if err != nil {
		t.Fatalf("creating handler: %v", err)
	}

	recorder := httptest.NewRecorder()
	responseHandler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	var response errorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if response.Error.Message != defaultErrorMessage {
		t.Fatalf("error message = %q, want %q", response.Error.Message, defaultErrorMessage)
	}
}
