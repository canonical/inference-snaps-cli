/*
Usage help for snapctl unset subcommand:

	snapctl [OPTIONS] unset [ConfKeys...]

	The unset command removes the provided configuration options as requested.

	$ snapctl unset name address

	All configuration changes are persisted at once, and only after the
	snap's configuration hook returns successfully.

	Nested values may be removed via a dotted path:

	$ snapctl unset user.name


	Help Options:
	-h, --help          Show this help message
*/

package snapctl

import (
	"fmt"
	"strings"
)

type unset struct {
	keys       []string
	validators []func() error
}

// Unset removes config options
// It takes one or more keys as input
func Unset(keys ...string) (cmd unset) {
	cmd.keys = keys

	cmd.validators = append(cmd.validators, func() error {
		for _, key := range cmd.keys {
			if strings.Contains(key, " ") {
				return fmt.Errorf("key must not contain spaces. Got: '%s'", key)
			}
		}
		return nil
	})

	return cmd
}

// Run executes the get command
func (cmd unset) Run() error {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return err
		}
	}

	// construct the command args
	// unset [ConfKeys...]
	var args []string
	// keys
	args = append(args, cmd.keys...)

	_, err := run("unset", args...)
	return err
}
