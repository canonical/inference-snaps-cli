package commands

import (
	"fmt"
	"strings"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

type serveUICommand struct {
	*common.Context

	// flags
	port         int
	host         string // bind address
	htmlDir      string
	capabilities string
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
	cobraCmd.Flags().IntVar(&cmd.port, "port", 8081, "HTTP bind port")
	cobraCmd.Flags().StringVar(&cmd.host, "host", "localhost", "HTTP bind address")
	cobraCmd.Flags().StringVar(&cmd.htmlDir, "dir", "./webui", "Directory to serve HTML files from")
	cobraCmd.Flags().StringVar(&cmd.capabilities, "capabilities", "text",
		fmt.Sprintf("Comma-separated list of capabilities (%s)",
			strings.Join(common.SupportedUICapabilities(), ", ")))

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

	var capabilities []string
	for _, c := range strings.Split(cmd.capabilities, ",") {
		c = strings.TrimSpace(c)
		switch c {
		case common.UICapabilityText, common.UICapabilityVision:
			capabilities = append(capabilities, c)
		default:
			return fmt.Errorf("unsupported capability: %s", c)
		}
	}

	config := common.WebUIConfig{
		OpenAIBaseURL: baseURL,
		Capabilities:  capabilities,
		InstanceName:  env.SnapInstanceName(),
		EngineName:    activeEngineName,
	}

	if cmd.Verbose {
		fmt.Println("OpenAI Base URL:", config.OpenAIBaseURL)
		fmt.Println("Capabilities:", config.Capabilities)
		fmt.Println("Instance name:", config.InstanceName)
		fmt.Println("Engine Name:", config.EngineName)
	}

	err = common.ServeUI(config, cmd.htmlDir, cmd.port, cmd.host)
	if err != nil {
		return fmt.Errorf("error serving UI: %v", err)
	}

	return nil
}
