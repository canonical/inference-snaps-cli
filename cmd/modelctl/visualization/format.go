package visualization

import (
	"bytes"
	"encoding/json"
	"fmt"

	"go.yaml.in/yaml/v4"
)

// Format is an output serialization format for machine information.
type Format string

const (
	// FormatYAML is a human-readable YAML output.
	FormatYAML Format = "yaml"
	// FormatJSON is a machine-readable, indented JSON output with kebab-cased keys.
	FormatJSON Format = "json"
)

// Marshal serializes the given visualization machine information using the
// requested format. It returns an error if the format is not recognized or if
// info is nil.
func Marshal(info *MachineInfo, f Format) ([]byte, error) {
	if info == nil {
		return nil, fmt.Errorf("cannot marshal nil machine info")
	}
	switch f {
	case FormatJSON:
		out, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("json: %w", err)
		}
		return append(out, '\n'), nil
	case FormatYAML:
		return marshalYAML(info)
	default:
		return nil, fmt.Errorf("unknown format %q", f)
	}
}

// marshalYAML renders machine information as human-readable YAML using a
// two-space indent.
func marshalYAML(info *MachineInfo) ([]byte, error) {
	var b bytes.Buffer
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(info); err != nil {
		return nil, fmt.Errorf("yaml: %w", err)
	}
	if err := enc.Close(); err != nil {
		return nil, fmt.Errorf("yaml: %w", err)
	}
	return b.Bytes(), nil
}
