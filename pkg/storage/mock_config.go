package storage

import (
	"maps"
	"os"
	"strings"
)

type mockConfig struct {
	values map[string]any
}

func NewMockConfig() Config {
	configValues := make(map[string]any)

	return &mockConfig{values: configValues}
}

func (c *mockConfig) Set(key, value string, confType configType) error {
	scopedKey := string(confType) + "." + key
	c.values[scopedKey] = value
	return nil
}

func (c *mockConfig) SetDocument(key string, value any, confType configType) error {
	scopedKey := string(confType) + "." + key
	c.values[scopedKey] = value
	return nil
}

func (c *mockConfig) Get(key string) (map[string]any, error) {
	configs, err := c.loadConfigs()
	if err != nil {
		return nil, err
	}

	// Filter to needed keys
	for k := range configs {
		if k != key && !strings.HasPrefix(k, key+".") {
			delete(configs, k)
		}
	}

	return configs, nil
}

func (c *mockConfig) GetAll() (map[string]any, error) {
	return c.loadConfigs()
}

func (c *mockConfig) Unset(key string, confType configType) error {
	scopedKey := string(confType) + "." + key
	delete(c.values, scopedKey)
	return nil
}

// loadConfigs loads all configurations as a flattened map, after applying precedence rules.
func (c *mockConfig) loadConfigs() (map[string]any, error) {
	values := make(map[string]any)

	for fullKey, value := range c.values {
		parts := strings.SplitN(fullKey, ".", 2)
		if len(parts) != 2 {
			continue
		}

		confKey := parts[0]
		key := parts[1]

		confValues, found := values[confKey]
		if !found {
			values[confKey] = map[string]any{key: value}
			continue
		}

		confValuesMap, ok := confValues.(map[string]any)
		if !ok {
			continue
		}
		confValuesMap[key] = value
	}

	envValues, err := c.getDefaultEnvVars()
	if err != nil {
		return nil, err
	}

	// Merge env values with mock values, giving env values precedence
	maps.Copy(values, envValues)

	finalMap := make(map[string]any)
	for _, k := range confPrecedence {
		if v, found := values[k]; found {
			vMap, ok := v.(map[string]any)
			if !ok {
				continue
			}
			maps.Copy(finalMap, c.flattenMap(vMap))
		}
	}

	return finalMap, nil
}

func (c *mockConfig) getDefaultEnvVars() (map[string]any, error) {
	envConfig := make(map[string]any)

	for _, entry := range os.Environ() {
		envKey, envValue, found := strings.Cut(entry, "=")
		if !found {
			continue
		}

		if !strings.HasPrefix(envKey, "DEFAULT_") {
			continue
		}

		normalizedKey := strings.TrimPrefix(envKey, "DEFAULT_")
		normalizedKey = strings.ToLower(strings.ReplaceAll(normalizedKey, "_", "-"))
		envConfig[normalizedKey] = envValue
	}

	if len(envConfig) == 0 {
		return map[string]any{}, nil
	}

	return map[string]any{
		string(Environment): envConfig,
	}, nil
}

// flattenMap creates a single-level map with dot-separated keys.
func (c *mockConfig) flattenMap(input map[string]any) map[string]any {
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
