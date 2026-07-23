package runtimes

const (
	// HTTP configuration keys
	HttpHostConfKey       = "http.host"
	HttpPortConfKey       = "http.port"
	HttpUnixSocketConfKey = "http.unix-socket"
	// WebSocket configuration keys
	WebSocketHostConfKey       = "ws.host"
	WebSocketPortConfKey       = "ws.port"
	WebSocketUnixSocketConfKey = "ws.unix-socket"

	// OpenAI server key for the runtime manifest.
	// This is used to identify OpenAI-compatible servers.
	OpenAiServerType = "openai"
)
