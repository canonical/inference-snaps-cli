package commands

import (
	"fmt"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

type serveUICommand struct {
	*common.Context

	// flags
	port       int
	host       string // bind address
	htmlDir    string
	inputImage bool
}

func ServeUI(ctx *common.Context) *cobra.Command {
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
	cobraCmd.Flags().IntVar(&cmd.port, "port", 8080, "HTTP bind port")
	cobraCmd.Flags().StringVar(&cmd.host, "host", "localhost", "HTTP bind address")
	cobraCmd.Flags().StringVar(&cmd.htmlDir, "dir", "./webui", "Directory to serve HTML files from")
	cobraCmd.Flags().BoolVar(&cmd.inputImage, "input-image", false, "Enable image input in the Web UI")

	return cobraCmd
}

func (cmd *serveUICommand) serveUI(_ *cobra.Command, args []string) error {

	baseURL, err := common.OpenAIEndpoint(cmd.Context)
	if err != nil {
		return fmt.Errorf("error getting OpenAI base URL: %v", err)
	}

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("error getting active engine: %v", err)
	}
	if activeEngineName == "" {
		return fmt.Errorf("no engine is active")
	}

	instanceName := env.SnapInstanceName()
	// Capitalize first letter, e.g. gemma3 -> Gemma3
	if instanceName != "" {
		instanceName = string(instanceName[0]-32) + instanceName[1:]
	}
	title := fmt.Sprintf("%s Inference Snap", instanceName)

	config := common.WebUIConfig{
		OpenAIBaseURL:  baseURL,
		SupportsImages: cmd.inputImage, // TODO: take this from an engine config?
		UITitle:        title,
		EngineName:     activeEngineName,
	}

	if cmd.Verbose {
		fmt.Println("Title:", config.UITitle)
		fmt.Println("OpenAI Base URL:", config.OpenAIBaseURL)
		fmt.Println("Supports Images:", config.SupportsImages)
		fmt.Println("Engine Name:", config.EngineName)
	}

	err = common.ServeUI(config, cmd.htmlDir, cmd.port, cmd.host)
	if err != nil {
		return fmt.Errorf("error serving UI: %v", err)
	}

	return nil
}
