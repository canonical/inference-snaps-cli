/*
Usage help for snapctl get subcommand:

	snapctl [OPTIONS] get [get-OPTIONS] [:<plug|slot>] [<keys>...]

	The get command prints configuration options for the current snap.

	$ snapctl get username
	frank

	If multiple option names are provided, a document is returned:

	$ snapctl get username password
	{
	"username": "frank",
	"password": "..."
	}

	Nested values may be retrieved via a dotted path:

	$ snapctl get author.name
	frank

	Values of interface connection settings may be printed with:

	$ snapctl get :myplug usb-vendor
	$ snapctl get :myslot path

	This will return the named setting from the local interface endpoint, whether a
	plug
	or a slot. Returning the setting from the connected snap's endpoint is also
	possible
	by explicitly requesting that via the --plug and --slot command line options:

	$ snapctl get :myplug --slot usb-vendor

	This requests the "usb-vendor" setting from the slot that is connected to
	"myplug".


	Help Options:
	-h, --help              Show this help message

	[get command options]
			--slot          return attribute values from the slot side of the
							connection
			--plug          return attribute values from the plug side of the
							connection
		-d                  always return document, even with single key
		-t                  strict typing with nulls and quoted strings

	[get command arguments]
	<keys>:                 option keys
*/

package snapctl

import (
	"fmt"
	"strings"
)

type get struct {
	options    []string
	_interface string
	keys       []string
	validators []func() error
}

// Get reads config options or interface attributes
// It takes an arbitrary number of keys as input
// It returns an object for setting the CLI arguments before running the command
func Get(key ...string) (cmd get) {
	cmd.keys = key

	cmd.validators = append(cmd.validators, func() error {
		for _, key := range cmd.keys {
			if strings.Contains(key, " ") {
				return fmt.Errorf("key must not contain spaces. Got: '%s'", key)
			}
		}
		if strings.HasPrefix(cmd._interface, ":") {
			return fmt.Errorf("interface plug/slot name must not contain colon as prefix")
		}
		return nil
	})

	return cmd
}

// Interface takes the plug or slot name
func (cmd get) Interface(name string) get {
	cmd._interface = name
	return cmd
}

// Document sets the following command option:
// -d	always return document, even with single key
func (cmd get) Document() get {
	cmd.options = append(cmd.options, "-d")
	return cmd
}

// Strict sets the following command option:
// -t	strict typing with nulls and quoted strings
func (cmd get) Strict() get {
	cmd.options = append(cmd.options, "-t")
	return cmd
}

// Run executes the get command
func (cmd get) Run() (string, error) {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return "", err
		}
	}

	// construct the command args
	// get [get-OPTIONS] [:<plug|slot>] [<keys>...]
	var args []string
	// options
	args = append(args, cmd.options...)
	// plug|slot, pre-validated
	if cmd._interface != "" {
		args = append(args, ":"+cmd._interface)
	}
	// keys
	args = append(args, cmd.keys...)

	return run("get", args...)
}
