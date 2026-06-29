package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/spf13/cobra"
)

type exportConfigCommand struct {
	*common.Context
}

func ExportConfig(ctx *common.Context) *cobra.Command {
	var cmd exportConfigCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "export-config",
		Short:             "Export the current configuration",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobraCmd
}

type sharedStatus struct {
	Engine    string            `json:"engine" yaml:"engine"`
	Services  map[string]string `json:"services" yaml:"services"`
	Endpoints map[string]string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Model     map[string]string `json:"model,omitempty" yaml:"model,omitempty"`
}

// sharedOpenai is to be deprecated in favor of sharedStatus.
// The Open WebUI snap first needs to be updated to consume sharedStatus
type sharedOpenai struct {
	BaseUrl string `json:"base_url"`
}

// run exports the inference snap status to a file where the content sharing interface can share it
// By default the status is written to $SNAP_COMMON/share/status.json
// and the openai configuration for open-webui is written to $SNAP_COMMON/share/openai.json
// The directory where these two files are written can be overridden by the $STATUS_SHARE_DIR environment variable
func (cmd *exportConfigCommand) run(_ *cobra.Command, _ []string) error {

	statusStr, err := common.StatusStruct(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting status: %v", err)
	}

	// Decouple internal status definition from shared one
	sharedStatusStr := &sharedStatus{
		Engine:    statusStr.Engine,
		Services:  statusStr.Services,
		Endpoints: statusStr.Endpoints,
		Model:     statusStr.Model,
	}

	statusJson, err := json.Marshal(sharedStatusStr)
	if err != nil {
		return fmt.Errorf("json: %v", err)
	}

	shareDir := os.Getenv("STATUS_SHARE_DIR")
	if shareDir == "" {
		shareDir = filepath.Join(os.Getenv("SNAP_COMMON"), "share")
	}

	if err := os.MkdirAll(shareDir, 0o755); err != nil {
		return fmt.Errorf("creating status share directory: %v", err)
	}

	statusFilePath := filepath.Join(shareDir, "status.json")
	if err := os.WriteFile(statusFilePath, statusJson, 0o644); err != nil {
		return fmt.Errorf("writing status.json: %v", err)
	}

	if endpoint, ok := statusStr.Endpoints["openai"]; ok {
		myOpenaiConfig := &sharedOpenai{
			BaseUrl: endpoint,
		}
		openaiJson, err := json.Marshal(myOpenaiConfig)
		if err != nil {
			return fmt.Errorf("json: %v", err)
		}
		openaiFilePath := filepath.Join(shareDir, "openai.json")
		if err := os.WriteFile(openaiFilePath, openaiJson, 0o644); err != nil {
			return fmt.Errorf("writing openai.json: %v", err)
		}
	}

	return nil
}
