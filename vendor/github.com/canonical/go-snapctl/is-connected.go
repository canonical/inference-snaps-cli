/*
Usage:
snapctl [OPTIONS] is-connected [is-connected-OPTIONS] <plug|slot>

	The is-connected command returns success if the given plug or slot of the
	calling snap is connected, and failure otherwise.

	$ snapctl is-connected plug
	$ echo $?
	1

	Snaps can only query their own plugs and slots - snap name is implicit and
	implied by the snapctl execution context.

	The --pid and --aparmor-label options can be used to determine whether
	a plug or slot is connected to the snap identified by the given
	process ID or AppArmor label.  In this mode, additional failure exit
	codes may be returned: 10 if the other snap is not connected but uses
	classic confinement, or 11 if the other process is not snap confined.

	The --pid and --apparmor-label options may only be used with slots of
	interface type "pulseaudio", "audio-record", or "cups-control".


	Help Options:
	-h, --help                Show this help message

	[is-connected command options]
			--pid=            Process ID for a plausibly connected process
			--apparmor-label= AppArmor label for a plausibly connected process

*/

package snapctl

import (
	"fmt"
	"strings"
)

type isConnected struct {
	plug       string
	validators []func() error
}

// IsConnected checks the connection status of a plug or slot
// It returns an object for setting the CLI arguments before running the command
func IsConnected(plug string) (cmd isConnected) {
	cmd.plug = plug

	cmd.validators = append(cmd.validators, func() error {
		if strings.Contains(plug, " ") {
			return fmt.Errorf("plug must not contain spaces. Got: '%s'", plug)
		}

		return nil
	})

	return cmd
}

// Run executes the get command
func (cmd isConnected) Run() (bool, error) {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return false, err
		}
	}

	// construct the command args
	// snapctl [OPTIONS] is-connected [is-connected-OPTIONS] <plug|slot>
	var args []string

	// plug
	args = append(args, cmd.plug)

	out, err := run("is-connected", args...)
	if err != nil && out == "" {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
