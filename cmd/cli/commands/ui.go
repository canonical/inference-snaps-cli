package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

type uiCommand struct {
	*common.Context
}

func Ui(ctx *common.Context) *cobra.Command {
	var cmd uiCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "ui",
		Short:             "Launch web UI",
		Long:              "Open the snap's builtin web user interface in the default browser",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobraCmd
}

func (cmd *uiCommand) run(_ *cobra.Command, _ []string) error {
	// Having all the components installed is not required, but it provides good feedback
	if err := common.WaitForComponents(cmd.Context); err != nil {
		return fmt.Errorf("waiting for component: %s", err)
	}

	// Get web ui url
	url, err := common.UiServerHttpUrl(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting ui server url: %s", err)
	}

	services, err := common.ServiceStatuses()
	if err != nil {
		return fmt.Errorf("getting service statuses: %v", err)
	}

	// Check ui server and engine server
	checkServices := []string{"server", "server-webui"}
	for _, service := range checkServices {
		uiServerStatus, ok := services[service]
		if !ok {
			return fmt.Errorf("%s: service not found", service)
		}
		if uiServerStatus == "inactive" {
			return fmt.Errorf("%s.%s not active\n\n%s",
				env.SnapInstanceName(), service, common.SuggestStartServer())
		}
	}

	// Wait until the openai server endpoint is ready to accept chat prompts. This should be handled in the webui in the future.
	chatBaseUrl, err := common.OpenAiEndpoint(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting OpenAI base URL: %v", err)
	}

	chatClient := common.ChatClient(chatBaseUrl, "", cmd.Verbose)
	err = chatClient.WaitChatServerReady() // prints same spinner as used for chat
	if err != nil {
		return err
	}

	// Print url and ask confirmation before opening
	fmt.Printf("Press Enter to open %s in the default browser ...\n", url)
	reader := bufio.NewReader(os.Stdin)
	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("waiting for Enter: %v", err)
	}

	// Use desktop portal to open URL in default browser
	err = exec.Command("xdg-open", url).Start()
	if err != nil {
		return fmt.Errorf("xdg-open: %v", err)
	}

	return nil
}
