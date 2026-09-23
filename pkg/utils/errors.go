package utils

import "errors"

// Error types that can be checked higher up in the caller chain
var (
	ErrInsufficientDiskSpaceForModel = errors.New("insufficient disk space for the selected model")
)
