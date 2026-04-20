package snap

import "github.com/canonical/go-snapctl/env"

func (*snap) InstanceName() string {
	return env.SnapInstanceName()
}

func InstanceName() string {
	return self.InstanceName()
}
