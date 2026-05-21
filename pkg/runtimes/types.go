package runtimes

type Manifest struct {
	Servers        map[string]Server `json:"servers" yaml:"servers"`
	Environment    []string          `json:"environment" yaml:"environment"`
	Layout         map[string]Layout `yaml:"layout"`
	expandedLayout map[string]Layout
	Components     []string `json:"components" yaml:"components"`
}

type Server struct {
	Protocol    string `json:"protocol" yaml:"protocol"`
	BasePath    string `json:"base-path" yaml:"base-path"`
	ConfigGroup string `json:"config-group" yaml:"config-group"`
}

type Layout struct {
	Symlink string `yaml:"symlink"`
}
