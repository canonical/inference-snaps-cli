package common

import "errors"

var (
	ErrPermissionDenied = errors.New("permission denied, try again with sudo")
	ErrGetActiveEngine  = errors.New("looking up active engine")
	ErrNoActiveEngine   = errors.New("no active engine")
)
