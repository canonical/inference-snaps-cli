package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebUIConfig struct {
	OpenAIBaseURL  string `json:"openAIBaseURL"`
	SupportsImages bool   `json:"supportsImages"`
	UITitle        string `json:"uiTitle"`
	EngineName     string `json:"engineName"`
}

func ServeUI(conf WebUIConfig, htmlDir string, port int, bindAddress string) error {

	mux := http.NewServeMux()

	// Serve configuration
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
