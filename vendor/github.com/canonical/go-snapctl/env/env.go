/*
 * Copyright (C) 2025 Canonical Ltd
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

package env

import (
	"os"
)

const (
	snapEnv         = "SNAP"
	snapCommonEnv   = "SNAP_COMMON"
	snapDataEnv     = "SNAP_DATA"
	snapInstNameEnv = "SNAP_INSTANCE_NAME"
	snapNameEnv     = "SNAP_NAME"
	snapRevEnv      = "SNAP_REVISION"
	snapVersionEnv  = "SNAP_VERSION"
)

// Getter functions for SNAP environment variables.

func Snap() string {
	return os.Getenv(snapEnv)
}

func SnapCommon() string {
	return os.Getenv(snapCommonEnv)
}

func SnapData() string {
	return os.Getenv(snapDataEnv)
}

func SnapInstanceName() string {
	return os.Getenv(snapInstNameEnv)
}

func SnapName() string {
	return os.Getenv(snapNameEnv)
}

func SnapRevision() string {
	return os.Getenv(snapRevEnv)
}

func SnapVersion() string {
	return os.Getenv(snapVersionEnv)
}
