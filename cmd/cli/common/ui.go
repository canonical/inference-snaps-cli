package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebUIConfig struct {
	OpenAIBaseURL string   `json:"openAIBaseURL"`
	Capabilities  []string `json:"capabilities"`
	InstanceName  string   `json:"instanceName"`
	EngineName    string   `json:"engineName"`
}

const (
	UICapabilityText   string = "text"
	UICapabilityVision string = "vision"
)

func SupportedUICapabilities() []string {
	return []string{UICapabilityText, UICapabilityVision}
}

func ServeUI(conf WebUIConfig, htmlDir string, port int, bindAddress string) error {

	mux := http.NewServeMux()

	// Serve configuration for the frontend
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(conf); err != nil {
			http.Error(w, "failed to encode config", http.StatusInternalServerError)
		}
	})

	// Serve the frontend static files
	mux.Handle("/", http.FileServer(http.Dir(htmlDir)))

	fmt.Printf("Serving %q on http://localhost:%d\n", htmlDir, port)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", bindAddress, port), mux)
}
