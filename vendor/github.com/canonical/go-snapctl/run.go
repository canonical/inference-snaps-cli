package snapctl

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/canonical/go-snapctl/log"
)

func run(subcommand string, subargs ...string) (string, error) {
	args := []string{subcommand}
	args = append(args, subargs...)

	log.Debugf("Executing 'snapctl %s'\n", strings.Join(args, " "))

	output, err := exec.Command("snapctl", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, output)
	}

	return strings.TrimSpace(string(output)), nil
}
