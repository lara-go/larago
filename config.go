package larago

import (
	"fmt"

	dotaccess "github.com/maxwellhealth/go-dotaccess"
)

// ConfigRepository used to store and get access to application config vars.
type ConfigRepository struct {
	config Config
}

// Env returns current environment name.
func (c *ConfigRepository) Env() string {
	return c.config.Env()
}

// Debug returs debug mode state.
func (c *ConfigRepository) Debug() bool {
	return c.config.Debug()
}

// Get value from config using dot-notation.
func (c *ConfigRepository) Get(key string) interface{} {
	value, err := dotaccess.Get(c.config, key)

	if err != nil {
		panic(fmt.Sprintf("Can not resolve config value: %s", key))
	}

	return value
}

// Set value to config using dot-notation.
func (c *ConfigRepository) Set(key string, value interface{}) {
	err := dotaccess.Set(c.config, key, value)

	if err != nil {
		panic(fmt.Sprintf("Can not resolve config value: %s", key))
	}
}
