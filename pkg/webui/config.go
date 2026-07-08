package webui

import (
	"fmt"
	"net/url"
)

type Config struct {
	OpenAIBaseURL string   `json:"openAIBaseURL"`
	Capabilities  []string `json:"capabilities"`
	InstanceName  string   `json:"instanceName"`
	EngineName    string   `json:"engineName"`
}

const (
	capabilityText     string = "text"
	capabilityVision   string = "vision"
	capabilityTools    string = "tools"
	capabilityThinking string = "thinking"
)

func SupportedCapabilities() []string {
	return []string{capabilityText, capabilityVision, capabilityTools, capabilityThinking}
}

func (c Config) Validate() error {

	// Validate OpenAI base URL
	if _, err := url.Parse(c.OpenAIBaseURL); err != nil {
		return fmt.Errorf("invalid OpenAI base URL: %w", err)
	}

	// Capabilities are forwarded as-is; the frontend ignores any it does not
	// recognize, so unknown capabilities are not rejected here.

	return nil
}
