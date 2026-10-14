/*
Usage help for snapctl stop subcommand:

	snapctl [OPTIONS] stop [stop-OPTIONS] <service>...

	The stop command stops the given services of the snap. If executed from the
	"configure" hook, the services will be stopped after the hook finishes.

	Help Options:
	-h, --help           Show this help message

	[stop command options]
			--disable    Disable the specified services (see man systemctl for
						details)
*/

package snapctl

import (
	"fmt"
	"strings"
)

type stop struct {
	services   []string
	options    []string
	validators []func() error
}

// Stop stops the services of the snap
// It takes an arbitrary number of service names as input
// It returns an object for setting the CLI arguments before running the command
func Stop(service ...string) (cmd stop) {
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

// Disable sets the following command option:
// --disable    Disable the specified services
func (cmd stop) Disable() stop {
	cmd.options = append(cmd.options, "--disable")
	return cmd
}

// Run executes the stop command
func (cmd stop) Run() error {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return err
		}
	}

	// construct the command args
	// stop [stop-OPTIONS] <service>...
	var args []string
	// options
	args = append(args, cmd.options...)
	// services
	args = append(args, cmd.services...)

	_, err := run("stop", args...)
	return err
}
