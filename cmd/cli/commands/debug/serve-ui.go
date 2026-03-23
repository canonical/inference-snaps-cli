package debug

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

type serveUiCommand struct {
	*common.Context

	// flags
	baseUrl string
	port    int
	htmlDir string
}

func ServeUICommand(ctx *common.Context) *cobra.Command {
	var cmd serveUiCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "serve-ui",
		Short:             "Run a debug web server to serve the Web UI",
		Hidden:            true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.serveUI,
	}

	// flags
	cobraCmd.Flags().StringVar(&cmd.baseUrl, "base-url", "http://localhost:8080/v1", "Base URL of the OpenAI-compatible server")
	cobraCmd.Flags().IntVar(&cmd.port, "port", 8080, "HTTP bind port")
	cobraCmd.Flags().StringVar(&cmd.htmlDir, "dir", "./webui", "Directory to serve HTML files from")

	return cobraCmd
}

type webUiConfig struct {
	OpenAIBaseURL  string `json:"openAIBaseURL"`
	SupportsImages bool   `json:"supportsImages"`
	UITitle        string `json:"uiTitle"`
	EngineName     string `json:"engineName"`
}

func (cmd *serveUiCommand) serveUI(_ *cobra.Command, args []string) error {

	config := webUiConfig{
		OpenAIBaseURL:  cmd.baseUrl,
		SupportsImages: true,
		UITitle:        "Debug Inference Snaps Web UI",
		EngineName:     "unset",
	}

	j, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Config: %s\n", j)

	mux := http.NewServeMux()

	// Serve configuration
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(config); err != nil {
			http.Error(w, "failed to encode config", http.StatusInternalServerError)
		}
	})

	// Serve the frontend static files
	mux.Handle("/", http.FileServer(http.Dir(cmd.htmlDir)))

	addr := fmt.Sprintf(":%d", cmd.port)
	fmt.Printf("Serving %q on http://localhost%s\n", cmd.htmlDir, addr)

	return http.ListenAndServe(addr, mux)
}
