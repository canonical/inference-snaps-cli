package debug

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

type serveUICommand struct {
	*common.Context

	// flags
	baseUrl string
	port    int
	htmlDir string
}

func ServeUICommand(ctx *common.Context) *cobra.Command {
	var cmd serveUICommand
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

func (cmd *serveUICommand) serveUI(_ *cobra.Command, args []string) error {

	config := common.WebUIConfig{
		OpenAIBaseURL:  cmd.baseUrl,
		SupportsImages: true,
		UITitle:        "Debug Inference Snaps Web UI",
		EngineName:     "unset",
	}

	j, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Config: %s\n", j)

	err := common.ServeUI(config, cmd.htmlDir, cmd.port, "localhost")
	if err != nil {
		return fmt.Errorf("error serving UI: %v", err)
	}

	return nil
}
