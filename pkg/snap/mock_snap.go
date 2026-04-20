package snap

import "fmt"

type mockSnap struct{}

func Mock() Snap {
	return &mockSnap{}
}

func (c *mockSnap) Restart() error {
	fmt.Println("[mock] Restarting snap")
	return nil
}

func (c *mockSnap) InstanceName() string {
	return "mock-snap"
}
