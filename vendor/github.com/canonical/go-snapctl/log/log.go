/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package log

import (
	"fmt"
	"log/syslog"
	"os"
)

var (
	debug            bool           // debug mode, set by snap option
	snapInstanceName string         // used as default syslog tag and tag prefix
	tag              string         // syslog tag and stderr prefix
	slog             *syslog.Writer // global syslog writer
)

// SetComponentName adds a component name to syslog tag as "snap.<snap-instance-name>.<component>"
// The default tag is just "snap.<snap-instance-name>".
// This function is NOT thread-safe. It should not be called concurrently with
// the other logging functions of this package.
func SetComponentName(component string) {
	// update global value
	tag = "snap." + snapInstanceName + "." + component
	Debugf("Changing syslog tag to: %s", tag)

	if err := setupSyslogWriter(tag); err != nil {
		Errorf("Error changing syslog tag: %s", err)
	}
}

// Debug writes the given input to syslog (sev=LOG_DEBUG) if snap `debug`
// configuration option is set to `true`.
// It formats similar to fmt.Sprint
func Debug(a ...interface{}) {
	if slog == nil {
		initialize()
	}

	if debug {
		slog.Debug(fmt.Sprint(a...))
	}
}

// Debugf writes the given input to syslog (sev=LOG_DEBUG) if snap `debug`
// configuration option is set to `true`.
// It formats similar to fmt.Sprintf
func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}

// Error writes the given input to syslog (sev=LOG_ERROR).
// It formats similar to fmt.Sprint
func Error(a ...interface{}) {
	if slog == nil {
		initialize()
	}

	msg := fmt.Sprint(a...)
	slog.Err(msg)
	// print to stderr as well so that snap command prints them on non-zero exit
	stderr(a...)
}

// Errorf writes the given input to syslog (sev=LOG_ERROR).
// It formats similar to fmt.Sprintf
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

// Fatal calls Error followed by os.Exit(1).
// It formats similar to fmt.Sprint
func Fatal(a ...interface{}) {
	Error(a...)
	os.Exit(1)
}

// Fatalf calls Errorf followed by os.Exit(1).
// It formats similar to fmt.Sprintf
func Fatalf(format string, a ...interface{}) {
	Errorf(format, a...)
	os.Exit(1)
}

// Info writes the given input to syslog (sev=LOG_INFO).
// It formats similar to fmt.Sprint
func Info(a ...interface{}) {
	if slog == nil {
		initialize()
	}

	slog.Info(fmt.Sprint(a...))
}

// Infof writes the given input to syslog (sev=LOG_INFO).
// It formats similar to fmt.Sprintf
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}

// Warn writes the given input to syslog (sev=LOG_WARNING).
// It formats similar to fmt.Sprint
func Warn(a ...interface{}) {
	if slog == nil {
		initialize()
	}

	slog.Warning(fmt.Sprint(a...))
}

// Warnf writes the given input to syslog (sev=LOG_WARNING).
// It formats similar to fmt.Sprintf
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}

// stderr writes the given input to standard error.
// It formats similar to fmt.Sprint, adds the global tag as prefix, and appends a newline
func stderr(a ...interface{}) {
	// Standard errors get collected with "snapd" as syslog app.
	// We add the tag as prefix to distinguish these from other snapd logs.
	fmt.Fprintf(os.Stderr, tag+": %s\n", a...)
}

func setupSyslogWriter(tag string) error {
	writer, err := syslog.New(syslog.LOG_INFO, tag)
	if err != nil {
		return err
	}
	// switch to new global writer only if no error
	slog = writer
	return nil
}
