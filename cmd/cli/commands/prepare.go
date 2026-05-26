package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	snapctl "github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type prepareCommand struct {
	*common.Context

	postRefresh bool
	install     bool
	noRestart   bool
	assumeYes   bool
}

func Prepare(ctx *common.Context) *cobra.Command {
	var cmd prepareCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare snap internal configurations",
		Long:  "Prepare internal configurations for the snap from a YAML file and perform necessary setup to run the snap.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  cmd.run,
	}

	cobraCmd.Flags().BoolVar(&cmd.postRefresh, "post-refresh", false, "triggered on post-refresh hook")
	cobraCmd.Flags().BoolVar(&cmd.install, "install", true, "triggered on install hook")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after setting the configuration")
	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for any prompt")

	return cobraCmd
}

func (cmd *prepareCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	if cmd.postRefresh || cmd.install {
		path := filepath.Join(env.Snap(), "var", "configurations", "configurations.yaml")
		if len(args) == 1 {
			path = args[0]
		}

		// no restart set to true to avoid double prompting. Prompting is handled in the caller function
		useEngineCmd := useEngineCommand{Context: cmd.Context, noRestart: true}

		var scoredEngines []engines.ScoredManifest
		connected, err := snapctl.IsConnected("hardware-observe").Run()
		if err != nil {
			return fmt.Errorf("checking hardware-observe connection: %v", err)
		}
		if connected {
			scoredEngines, err = common.ScoreEnginesWithSpinner(cmd.Context)
			if err != nil {
				return fmt.Errorf("scoring engines: %v", err)
			}
		} else {
			fmt.Fprintln(os.Stderr, "hardware-observe interface not auto connected. Skip auto engine selection.")
		}

		return cmd.load(path, cmd.Context, scoredEngines, &useEngineCmd)
	} else {
		return fmt.Errorf("prepare command should be run with at least one of --post-refresh or --install flag")
	}
}

func (cmd *prepareCommand) load(path string, ctx *common.Context, scoredEngines []engines.ScoredManifest, useEngineCmd *useEngineCommand) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading configurations from %q: %w", path, err)
	}

	var conf map[string]any
	if err := yaml.Unmarshal(content, &conf); err != nil {
		return fmt.Errorf("parsing configurations from %q: %w", path, err)
	}

	changed, err := cmd.setConfigurations(conf)
	if err != nil {
		return err
	}

	if err := cmd.createSymbolicLinksForWSL(); err != nil {
		return fmt.Errorf("creating symbolic links for WSL: %w", err)
	}

	if scoredEngines != nil {
		if err := useEngineCmd.autoSelectScoredEngine(scoredEngines); err != nil {
			return fmt.Errorf("auto-selecting engine: %w", err)
		}
	}

	if !cmd.noRestart && changed {
		return common.PromptRestartToApplyChanges(ctx, cmd.assumeYes)
	}

	return nil
}

func (cmd *prepareCommand) createSymbolicLinksForWSL() error {
	links := []struct{ target, link string }{
		{
			target: "/var/lib/snapd/hostfs/usr/lib/wsl/lib",
			link:   filepath.Join(env.SnapCommon(), "usr", "lib", "wsl", "lib"),
		},
		{
			target: "/var/lib/snapd/hostfs/usr/lib/wsl/drivers",
			link:   filepath.Join(env.SnapCommon(), "usr", "lib", "wsl", "drivers"),
		},
	}

	for _, l := range links {
		if err := utils.CreateSymlink(l.target, l.link); err != nil {
			return fmt.Errorf("creating symlink %q -> %q: %w", l.link, l.target, err)
		}
	}

	return nil
}

func (cmd *prepareCommand) setConfigurations(conf map[string]any) (bool, error) {

	// cannot check for errors because if the config is not set, it will be set to nil and the error will be ignored
	before, _ := cmd.Config.GetAll()

	for configurationType := range conf {
		switch configurationType {
		case "package":
			for k, v := range conf[configurationType].(map[string]any) {
				if err := cmd.Config.SetDocument(k, v, storage.PackageConfig); err != nil {
					return false, fmt.Errorf("setting package configurations: %w", err)
				}
			}
		case "user":
			for k, v := range conf[configurationType].(map[string]any) {
				if err := cmd.Config.SetDocument(k, v, storage.UserConfig); err != nil {
					return false, fmt.Errorf("setting user configurations: %w", err)
				}
			}
		case "engine":
			for k, v := range conf[configurationType].(map[string]any) {
				if err := cmd.Config.SetDocument(k, v, storage.EngineConfig); err != nil {
					return false, fmt.Errorf("setting engine configurations: %w", err)
				}
			}
		default:
			return false, fmt.Errorf("unknown configuration type %q", configurationType)
		}
	}

	after, err := cmd.Config.GetAll()
	if err != nil {
		return false, fmt.Errorf("getting updated configurations: %w", err)
	}

	return !reflect.DeepEqual(before, after), nil
}


