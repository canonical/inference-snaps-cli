package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/constants"
	"github.com/canonical/inference-snaps-cli/v2/pkg/fallbackserver"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/spf13/cobra"
)

type runCommand struct {
	*common.Context

	// flags
	waitForComponents bool
	shareProvider     string
	fallbackServer    bool
}

func Run(ctx *common.Context) *cobra.Command {
	var cmd runCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "run <command>",
		Short: "Run a subprocess",
		Long: "Run a command in the engine's environment\n\n" +
			"Use run to execute a program as a sub-process, within the active engine's environment.\n" +
			"To pass arguments to the program itself, separate the command and its arguments with\n" +
			"double dashes (--) from the run command and its flags. ",
		Example: "  modelctl run env\n" +
			"  modelctl run -- echo \"Hello World!\"\n" +
			"  modelctl run --share-provider -- python3 -m http.server",
		Hidden:            true,
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	// --share-provider [path]
	cobraCmd.Flags().StringVar(&cmd.shareProvider, "share-provider", "", "write provider env file to a shared directory")
	cobraCmd.Flags().Lookup("share-provider").NoOptDefVal = cmd.defaultProviderDirectoryPath()
	// --fallback-server
	cobraCmd.Flags().BoolVar(&cmd.fallbackServer, "fallback-server", false, "if the command fails to run, start a fallback server")

	return cobraCmd
}

func (cmd *runCommand) run(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("unexpected number of arguments, expected at least 1 got %d", len(args))
	}

	// Components are required for loading the engine environment
	if err := common.WaitForComponents(cmd.Context); err != nil {
		return fmt.Errorf("waiting for component: %s", err)
	}

	clean, err := common.LoadEngineEnvironment(cmd.Context)
	if errors.Is(err, common.ErrNoActiveModel) {
		return cmd.startFallbackServer(err)
	}
	if err != nil {
		return fmt.Errorf("loading engine environment: %v", err)
	}

	// NOTE: defer does not run on SIGTERM or SIGKILL. It only runs when the child process exits.
	// TODO: add signal handling to intercept SIGTERM and invoke clean() before exiting.
	defer clean()

	if err := cmd.processEnvConfigs(); err != nil {
		return fmt.Errorf("processing env configs: %v", err)
	}
	if err := cmd.writeShareProviderEnv(); err != nil {
		return fmt.Errorf("writing share provider env: %v", err)
	}

	command := args[0]

	execCmd := exec.Command(command, args[1:]...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	commandErr := execCmd.Run()

	// systemd normally sends SIGTERM to every process in the service's
	// control group. If only the child receives it, Run() returns an ExitError.
	// Treat that as an intentional stop rather than a failure.
	if commandErr != nil && !commandStopped(commandErr) {
		return cmd.startFallbackServer(commandErr)
	}

	return nil
}

func (cmd *runCommand) startFallbackServer(commandErr error) error {
	if cmd.fallbackServer {
		// For now only serve a fallback server if the engine defines an openai endpoint
		url, err := common.OpenAiBaseUrl(cmd.Context)
		if err != nil && errors.Is(err, common.ErrNoOpenAiServer) {
			return commandErr
		} else if err != nil {
			return fmt.Errorf("getting OpenAI base URL: %v", err)
		}

		var servedErrorMessages []string

		statusStr, err := common.SnapStatus(cmd.Context)
		if err != nil {
			if errors.Is(err, common.ErrNoActiveModel) {
				servedErrorMessages = append(servedErrorMessages, "No active model is set. Please set an active model and try again.")
			} else {
				return fmt.Errorf("getting status: %v", err)
			}
		}
		if statusStr == nil {
			return fmt.Errorf("empty status reported")
		}

		// Report all notices as errors to explain why the server failed
		servedErrorMessages = append(servedErrorMessages, statusStr.Notices...)

		fmt.Println("Starting fallback server...")

		if err := fallbackserver.Run(url, servedErrorMessages); err != nil {
			return fmt.Errorf("running fallback server: %v", err)
		}
	} else {
		return commandErr
	}
	return nil
}

func commandStopped(err error) bool {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}

	status, ok := exitErr.ProcessState.Sys().(syscall.WaitStatus)
	if !ok || !status.Signaled() {
		return false
	}

	signal := status.Signal()
	return signal == syscall.SIGTERM || signal == syscall.SIGINT
}

func (cmd *runCommand) processEnvConfigs() error {
	envConfigs, err := cmd.Config.Get("env")
	if err != nil {
		return fmt.Errorf("getting configs: %v", err)
	}

	envVars := make(map[string]any, len(envConfigs))
	for k, v := range envConfigs {
		// Convert env keys (my-key) to environment variable names (MY_KEY)
		name, ok := strings.CutPrefix(k, storage.EnvKeyPrefix)
		if !ok {
			return fmt.Errorf("unexpected config key %q: expected prefix %q", k, storage.EnvKeyPrefix)
		}
		name = strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
		envVars[name] = fmt.Sprintf("%v", v)
	}

	err = utils.SetEnvironmentVariables(envVars)
	if err != nil {
		return fmt.Errorf("setting environment variables: %v", err)
	}
	return nil
}

func (cmd *runCommand) defaultProviderDirectoryPath() string {
	return os.ExpandEnv(constants.DefaultShareProviderPath)
}

func (cmd *runCommand) writeShareProviderEnv() error {
	if cmd.shareProvider == "" {
		return nil
	}

	if err := os.MkdirAll(cmd.shareProvider, 0o755); err != nil {
		return fmt.Errorf("creating provider env directory: %v", err)
	}

	if err := cmd.Cache.SetSharedProviderDirectory(cmd.shareProvider); err != nil {
		return fmt.Errorf("saving shared provider directory: %v", err)
	}

	providerEnvPath := filepath.Join(cmd.shareProvider, "provider.env")
	var content strings.Builder

	content.WriteString("SNAP_NAME=")
	content.WriteString(snap.SnapName())
	content.WriteString("\n")

	content.WriteString("SNAP_INSTANCE_NAME=")
	content.WriteString(snap.InstanceName())
	content.WriteString("\n")

	baseURL, err := common.OpenAiBaseUrl(cmd.Context)
	if err != nil && !errors.Is(err, common.ErrNoOpenAiServer) && !errors.Is(err, common.ErrOpenAiServerNoUrl) {
		return fmt.Errorf("getting OpenAI base URL: %v", err)
	} else if err == nil {
		content.WriteString("OPENAI_BASE_URL=")
		content.WriteString(baseURL)
		content.WriteString("\n")
	}

	runtime, err := common.CurrentRuntimeManifest(cmd.Context)
	if err != nil && !errors.Is(err, common.ErrEngineNoRuntime) {
		return fmt.Errorf("getting current runtime manifest: %v", err)
	}
	for _, server := range runtime.Servers {
		if server.IsUnixProtocol() {
			if server.Namespace != "" {
				content.WriteString(strings.ReplaceAll(strings.ToUpper(server.Namespace), "-", "_"))
				content.WriteString("_")
			}
			content.WriteString("UNIX_SOCKET")
			content.WriteString("=")
			content.WriteString(server.UnixSocketName())
			content.WriteString("\n")
		}
	}

	tmpPath := providerEnvPath + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(content.String()), 0o644); err != nil {
		return fmt.Errorf("writing provider env file: %v", err)
	}
	if err := os.Rename(tmpPath, providerEnvPath); err != nil {
		return fmt.Errorf("renaming provider env file: %v", err)
	}

	return nil
}
