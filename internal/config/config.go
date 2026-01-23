package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Config holds all application configuration
type Config struct {
	General GeneralConfig `json:"general"`
	Engine  EngineConfig  `json:"engine"`
	Prompt  PromptConfig  `json:"prompt"`
}

// Default returns a Config with default values
func Default() *Config {
	return &Config{
		General: DefaultGeneralConfig(),
		Engine:  DefaultEngineConfig(),
		Prompt:  DefaultPromptConfig(),
	}
}

var (
	configDir  string
	configOnce sync.Once
)

// getConfigDir returns the configuration directory path
func getConfigDir() string {
	configOnce.Do(func() {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			// Fallback to home directory
			home, _ := os.UserHomeDir()
			configDir = filepath.Join(home, ".tons")
		} else {
			configDir = filepath.Join(userConfigDir, "tons")
		}
	})
	return configDir
}

// configPath returns the full path to the config file
func configPath() string {
	return filepath.Join(getConfigDir(), "config.json")
}

// Load reads the configuration from disk
// Returns default config if file doesn't exist
func Load() (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	dir := getConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath(), data, 0644)
}

// Reset restores the configuration to default values and saves
func (c *Config) Reset() error {
	*c = *Default()
	return c.Save()
}
