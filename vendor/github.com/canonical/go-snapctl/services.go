/*
Usage help for snapctl services subcommand:

	snapctl [OPTIONS] services [<service>...]

	The services command lists information about the services specified.
*/

package snapctl

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type services struct {
	names      []string
	validators []func() error
}

// service status object
type service struct {
	Startup string
	Current string
	Notes   string
}

// Services lists information about the services
// It takes zero or more service names as input
// It returns an object for setting the CLI arguments before running the command
func Services(name ...string) (cmd services) {
	cmd.names = name

	cmd.validators = append(cmd.validators, func() error {
		for _, name := range cmd.names {
			if strings.Contains(name, " ") {
				return fmt.Errorf("service name must not contain spaces. Got: '%s'", name)
			}
		}
		return nil
	})

	return cmd
}

// Run executes the services command
func (cmd services) Run() (map[string]service, error) {
	// validate all input
	for _, validate := range cmd.validators {
		if err := validate(); err != nil {
			return nil, err
		}
	}

	// construct the command args
	// services [<service>...]
	var args []string
	// service names
	args = append(args, cmd.names...)

	output, err := run("services", args...)
	if err != nil {
		return nil, err
	}

	return cmd.parseOutput(output)
}

func (cmd services) parseOutput(output string) (map[string]service, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))

	// throw away the header:
	// Service   Startup   Current   Notes
	scanner.Scan()

	services := make(map[string]service)
	for scanner.Scan() {
		line := scanner.Text()

		// Split by whitespaces up to four parts.
		// The last part is for notes which may contain spaces in itself.
		cells := regexp.MustCompile("[[:space:]]+").Split(line, 4)
		if len(cells) != 4 {
			return nil, fmt.Errorf("unexpected snapctl output: expected 4 columns, got: %d", len(cells))
		}

		serviceName := cells[0]
		startup := cells[1]
		current := cells[2]
		notes := cells[3]

		services[serviceName] = service{
			Startup: startup,
			Current: current,
			Notes:   notes,
		}
	}

	return services, nil
}
