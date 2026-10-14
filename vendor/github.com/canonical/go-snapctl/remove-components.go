/*
Usage:
  snapctl [OPTIONS] remove <snap|snap+comp|+comp>...

The remove command removes components.


Help Options:
  -h, --help                        Show this help message

[remove command arguments]
  <snap|snap+comp|+comp>:           Components to be removed (snap must be
                                    the caller snap if specified).
*/

package snapctl

import (
	"fmt"
	"strings"
)

type removeComponents struct {
	components []string
	validators []func() error
}

// RemoveComponents removes components of the snap
// It takes an arbitrary number of component names as input
// It returns an object for setting the CLI arguments before running the command
func RemoveComponents(components ...string) (cmd removeComponents) {
	cmd.components = append(cmd.components, components...)

	cmd.validators = append(cmd.validators, func() error {
		if len(cmd.components) == 0 {
			return fmt.Errorf("at least one component must be specified")
		}
		return nil
	})

	cmd.validators = append(cmd.validators, func() error {
		for _, component := range cmd.components {
			if strings.Contains(component, " ") {
				return fmt.Errorf("component names must not contain spaces. Got: '%s'", component)
			}
		}
		return nil
	})

	return cmd
}

// Run executes the remove command
func (cmd removeComponents) Run() error {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return err
		}
	}

	// construct the command args
	// remove [remove-OPTIONS] <snap|snap+comp|+comp>... - we only support <+comp>... for now
	var args []string

	for _, component := range cmd.components {
		args = append(args, "+"+component)
	}

	_, err := run("remove", args...)
	return err
}
