package snap

import (
	"github.com/canonical/go-snapctl"
)

// Restart restarts all or a subset of snap services.
// To restart all, run without arguments.
func (*snap) Restart(service ...string) error {
	return snapctl.Restart(service...).Run()
}

func Restart(service ...string) error {
	return self.Restart(service...)
}
