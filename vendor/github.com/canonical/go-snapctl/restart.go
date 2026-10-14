/*
Usage help for snapctl restart subcommand:

	snapctl [OPTIONS] restart [restart-OPTIONS] <service>...

	The restart command restarts the given services of the snap. If executed from
	the
	"configure" hook, the services will be restarted after the hook finishes.

	Help Options:
	-h, --help           Show this help message

	[restart command options]
			--reload     Reload the given services if they support it (see man
						systemctl for details)
*/

package snapctl

import (
	"fmt"
	"strings"
)

type restart struct {
	services   []string
	options    []string
	validators []func() error
}

// Restart restarts the services of the snap
// It takes an arbitrary number of service names as input
// It returns an object for setting the CLI arguments before running the command
func Restart(service ...string) (cmd restart) {
	cmd.services = service

	cmd.validators = append(cmd.validators, func() error {
		for _, key := range cmd.services {
			if strings.Contains(key, " ") {
				return fmt.Errorf("service names must not contain spaces. Got: '%s'", key)
			}
		}
		return nil
	})

	return cmd
}

// Reload sets the following command option:
// --reload     Reload the given services if they support it
func (cmd restart) Reload() restart {
	cmd.options = append(cmd.options, "--reload")
	return cmd
}

// Run executes the restart command
func (cmd restart) Run() error {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return err
		}
	}

	// construct the command args
	// restart [restart-OPTIONS] <service>...
	var args []string
	// options
	args = append(args, cmd.options...)
	// services
	args = append(args, cmd.services...)

	_, err := run("restart", args...)
	return err
}
