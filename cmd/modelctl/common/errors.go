package common

import "errors"

// Error types that can be checked higher up in the caller chain
var (
	ErrPermissionDenied = errors.New("permission denied, try again with sudo")
	ErrNoActiveEngine   = errors.New("no active engine")
	ErrNoActiveRuntime  = errors.New("no active runtime")
	ErrNoActiveModel    = errors.New("no active model")
	ErrNoOpenAiServer   = errors.New("no OpenAI server available")
)

// Strings that are commonly used in error chains, but should not be used as error types
const (
	LookingUpActiveEngine         = "looking up active engine"
	LookingUpActiveModel          = "looking up active model"
	LoadingEngineManifest         = "loading engine manifest"
	LoadingModelManifests         = "loading model manifests"
	InsufficientDiskSpaceForModel = "insufficient disk space for the selected model"
)
