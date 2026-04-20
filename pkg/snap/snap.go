package snap

type Snap interface {
	Restart() error
	InstanceName() string
}

type snap struct{}

var self snap
