package storage

import (
	"encoding/json"
	"errors"
	"maps"
	"strings"
)

type config struct {
	storage storage
}

func NewConfig() *config {
	return &config{
		storage: NewSnapctlStorage(), // hardcoded since that's the only supported backend
	}
}

const configKeyPrefix = "config"

type configType string

// config precedence, from lowest to highest
var confPrecedence = []configType{
	PackageConfig, // values set by the package
	EngineConfig,  // values set by the active engine, overriding package values
	UserConfig,    // values set by the user, overriding all others
}

// config types
const (
	PackageConfig configType = "package"
	EngineConfig  configType = "engine"
	UserConfig    configType = "user"
)

func (c config) Set(key, value string, confType configType) error {
	return c.storage.Set(c.nestKeys(confType, key), value)
}

// Deprecated
// Remove once migration from config to cache is complete
func (c config) SetDocument(key string, value any, confType configType) error {
	return c.storage.SetDocument(c.nestKeys(confType, key), value)
}

// Get returns one configuration field, after applying precedence rules
// If the value is an object, it is returned as a JSON string
// TODO: change this to return a flattened map[string]any to better support object values
func (c *config) Get(key string) (string, error) {
	var value string

	// Load configurations in the order of precedence
	for _, confType := range confPrecedence {
		data, err := c.storage.Get(c.nestKeys(confType, key))
		if err != nil {
			if errors.Is(err, ErrorNotFound) {
				continue
			}
			return "", err
		}
		value = string(data)
	}

	return value, nil
}

// GetAll returns all configurations as a flattened map, after applying precedence rules
func (c *config) GetAll() (map[string]any, error) {
	data, err := c.storage.Get(configKeyPrefix)
	if err != nil {
		return nil, err
	}
	var valMap map[string]any
	err = json.Unmarshal(data, &valMap)
	if err != nil {
		return nil, err
	}

	// Drop the "engines" object (manifests and scores) from output
	// TODO: remove once no longer using the deprecated "engines" config
	delete(values[string(PackageConfig)].(map[string]any), "engines")

	var finalMap = make(map[string]any)

	// Load configurations in the order of precedence
	for _, k := range confPrecedence {
		if v, found := values[string(k)]; found {
			maps.Copy(
				finalMap,
				c.flattenMap(v.(map[string]any)),
			)
		}
	}

	return finalMap, nil
}

func (c config) Unset(key string, confType configType) error {
	return c.storage.Unset(c.nestKeys(confType, key))
}

// flattenMap creates a single-level map with dot-separated keys
func (c *config) flattenMap(input map[string]any) map[string]any {
	flatMap := make(map[string]any)

	var recurse func(map[string]any, string)
	recurse = func(m map[string]any, prefix string) {
		for k, v := range m {
			fullKey := k
			if prefix != "" {
				fullKey = prefix + "." + k
			}
			switch val := v.(type) {
			case map[string]any:
				recurse(val, fullKey)
			default:
				flatMap[fullKey] = val
			}
		}
	}
	recurse(input, "")

	return flatMap
}

// nestKeys creates a dot-separated key with the expected prefix
func (c *config) nestKeys(confType configType, key string) string {
	if key == "." { // special case, referencing the parent
		return strings.Join([]string{configKeyPrefix, string(confType)}, ".")
	} else {
		return strings.Join([]string{configKeyPrefix, string(confType), key}, ".")
	}

}
