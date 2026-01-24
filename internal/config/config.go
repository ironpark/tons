package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Config holds all application configuration
type Config struct {
	mu      sync.RWMutex  `json:"-"`
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
	c.mu.RLock()
	defer c.mu.RUnlock()

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
	c.mu.Lock()
	defaultCfg := Default()
	c.General = defaultCfg.General
	c.Engine = defaultCfg.Engine
	c.Prompt = defaultCfg.Prompt
	c.mu.Unlock()

	return c.Save()
}

// Snapshot returns a deep copy of the current configuration
func (c *Config) Snapshot() *Config {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := &Config{
		General: c.General,
		Engine:  c.Engine,
		Prompt:  c.Prompt,
	}

	// Deep copy slices in TerminalAgentConfig
	if c.Engine.TerminalAgent.ClaudeCode.Args != nil {
		snapshot.Engine.TerminalAgent.ClaudeCode.Args = make([]string, len(c.Engine.TerminalAgent.ClaudeCode.Args))
		copy(snapshot.Engine.TerminalAgent.ClaudeCode.Args, c.Engine.TerminalAgent.ClaudeCode.Args)
	}
	if c.Engine.TerminalAgent.GeminiCLI.Args != nil {
		snapshot.Engine.TerminalAgent.GeminiCLI.Args = make([]string, len(c.Engine.TerminalAgent.GeminiCLI.Args))
		copy(snapshot.Engine.TerminalAgent.GeminiCLI.Args, c.Engine.TerminalAgent.GeminiCLI.Args)
	}
	if c.Engine.TerminalAgent.Codex.Args != nil {
		snapshot.Engine.TerminalAgent.Codex.Args = make([]string, len(c.Engine.TerminalAgent.Codex.Args))
		copy(snapshot.Engine.TerminalAgent.Codex.Args, c.Engine.TerminalAgent.Codex.Args)
	}

	return snapshot
}

// Restore applies the given snapshot to the current configuration
func (c *Config) Restore(snapshot *Config) {
	if snapshot == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.General = snapshot.General
	c.Engine = snapshot.Engine
	c.Prompt = snapshot.Prompt

	// Deep copy slices
	if snapshot.Engine.TerminalAgent.ClaudeCode.Args != nil {
		c.Engine.TerminalAgent.ClaudeCode.Args = make([]string, len(snapshot.Engine.TerminalAgent.ClaudeCode.Args))
		copy(c.Engine.TerminalAgent.ClaudeCode.Args, snapshot.Engine.TerminalAgent.ClaudeCode.Args)
	}
	if snapshot.Engine.TerminalAgent.GeminiCLI.Args != nil {
		c.Engine.TerminalAgent.GeminiCLI.Args = make([]string, len(snapshot.Engine.TerminalAgent.GeminiCLI.Args))
		copy(c.Engine.TerminalAgent.GeminiCLI.Args, snapshot.Engine.TerminalAgent.GeminiCLI.Args)
	}
	if snapshot.Engine.TerminalAgent.Codex.Args != nil {
		c.Engine.TerminalAgent.Codex.Args = make([]string, len(snapshot.Engine.TerminalAgent.Codex.Args))
		copy(c.Engine.TerminalAgent.Codex.Args, snapshot.Engine.TerminalAgent.Codex.Args)
	}
}
