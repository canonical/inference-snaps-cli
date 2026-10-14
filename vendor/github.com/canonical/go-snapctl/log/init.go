package log

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/canonical/go-snapctl/env"
)

func initialize() {
	if err := initializeLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing go-snapctl/log: %v\n", err)
		fmt.Fprintln(os.Stderr, "The initialization can only be done from the snap environment.")
		os.Exit(1)
	}
}

func initializeLogger() error {
	snapInstanceName = env.SnapInstanceName()
	if snapInstanceName == "" {
		return fmt.Errorf("SNAP_INSTANCE_NAME environment variable not set")
	}
	tag = "snap." + snapInstanceName

	value, err := exec.Command("snapctl", "get", "debug").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error getting value of debug snap option: %v", err)
	}
	debug = (string(bytes.TrimSpace(value)) == "true")

	if err := setupSyslogWriter(tag); err != nil {
		return fmt.Errorf("error setting up syslog writer: %v", err)
	}

	return nil
}
