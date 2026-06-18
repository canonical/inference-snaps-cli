package storage

import (
	"errors"
	"fmt"
)

func (c *config) Migrate() error {
	if err := c.migratePassthroughEnv(); err != nil {
		return fmt.Errorf("migrating passthrough environment variables: %v", err)
	}
	return nil
}

// migratePassthroughEnv migrates v1 passthrough environment variables
// passthrough.env.<var-name>=<value> to the v2 env.<var-name>=<value> structure.
func (c *config) migratePassthroughEnv() error {
	values, err := c.storage.Get(c.nestKeys(UserConfig, "passthrough.env"))
	if err != nil {
		if errors.Is(err, ErrorNotFound) {
			return nil
		}
		return err
	}

	for varName, varValue := range values {
		if err := c.Set("env."+varName, fmt.Sprint(varValue), UserConfig); err != nil {
			return fmt.Errorf("setting env.%s: %v", varName, err)
		}
		if err := c.Unset("passthrough.env."+varName, UserConfig); err != nil {
			return fmt.Errorf("unsetting passthrough.env.%s: %v", varName, err)
		}
	}

	return nil
}
