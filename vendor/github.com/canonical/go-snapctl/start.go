/*
Usage help for snapctl start subcommand:

	snapctl [OPTIONS] start [start-OPTIONS] <service>...

	The start command starts the given services of the snap. If executed from the
	"configure" hook, the services will be started after the hook finishes.

	Help Options:
	-h, --help           Show this help message

	[start command options]
			--enable     Enable the specified services (see man systemctl for
						details)
*/

package snapctl

import (
	"fmt"
	"strings"
)

type start struct {
	services   []string
	options    []string
	validators []func() error
}

// Start starts the services of the snap
// It takes an arbitrary number of service names as input
// It returns an object for setting the CLI arguments before running the command
func Start(service ...string) (cmd start) {
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

// Enable sets the following command option:
// --enable     Enable the specified services
func (cmd start) Enable() start {
	cmd.options = append(cmd.options, "--enable")
	return cmd
}

// Run executes the start command
func (cmd start) Run() error {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return err
		}
	}

	// construct the command args
	// start [start-OPTIONS] <service>...
	var args []string
	// options
	args = append(args, cmd.options...)
	// services
	args = append(args, cmd.services...)

	_, err := run("start", args...)
	return err
}
