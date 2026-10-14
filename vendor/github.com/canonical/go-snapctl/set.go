/*
Usage help for snapctl set subcommand:

	snapctl [OPTIONS] set [set-OPTIONS] [:<plug|slot>] [key=value...]

	The set command changes the provided configuration options as requested.

	$ snapctl set username=frank password=$PASSWORD

	All configuration changes are persisted at once, and only after the hook
	returns successfully.

	Nested values may be modified via a dotted path:

	$ snapctl set author.name=frank

	Configuration option may be unset with exclamation mark:
	$ snapctl set author!

	Plug and slot attributes may be set in the respective prepare and connect hooks
	by
	naming the respective plug or slot:

	$ snapctl set :myplug path=/dev/ttyS0

	Help Options:
	-h, --help              Show this help message

	[set command options]
		-s                  parse the value as a string
		-t                  parse the value strictly as JSON document
*/

package snapctl

import (
	"errors"
	"fmt"
	"strings"
)

type set struct {
	options    []string
	_interface string
	keyValues  []string
	validators []func() error
}

// Set writes config options or interface attributes
// It takes one or more key-value pairs as input
// It returns an object for setting optional CLI arguments before running the command
func Set(keyValue ...string) (cmd set) {
	cmd.keyValues = keyValue

	cmd.validators = append(cmd.validators, func() error {
		for i := 0; i < len(cmd.keyValues); i += 2 {
			key := cmd.keyValues[i]
			if strings.Contains(key, " ") {
				return fmt.Errorf("key must not contain spaces. Got: '%s'", key)
			}
		}
		if len(cmd.keyValues)%2 != 0 {
			return fmt.Errorf("key-value inputs must be even and in pairs, got %d",
				len(cmd.keyValues))
		}
		return nil
	})

	return cmd
}

// Interface takes the plug or slot name
func (cmd set) Interface(name string) set {
	cmd._interface = name

	cmd.validators = append(cmd.validators, func() error {
		if strings.HasPrefix(cmd._interface, ":") {
			return errors.New("interface plug/slot name must not contain colon as prefix")
		}
		return nil
	})

	return cmd
}

// Document sets the following command option:
// -t	parse the value strictly as JSON document
func (cmd set) Document() set {
	cmd.options = append(cmd.options, "-t")
	return cmd
}

// String sets the following command option:
// -s   parse the value as a string
func (cmd set) String() set {
	cmd.options = append(cmd.options, "-s")
	return cmd
}

// Run executes the get command
func (cmd set) Run() error {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return err
		}
	}

	// construct the command args
	// set [set-OPTIONS] [:<plug|slot>] [key=value...]
	var args []string
	// options
	args = append(args, cmd.options...)
	// plug|slot
	if cmd._interface != "" {
		args = append(args, ":"+cmd._interface)
	}
	// key-values
	for i := 0; i < len(cmd.keyValues); i += 2 {
		args = append(args, fmt.Sprintf("%s=%s",
			cmd.keyValues[i],
			cmd.keyValues[i+1]))
	}

	_, err := run("set", args...)
	return err
}
