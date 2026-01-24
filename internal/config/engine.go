package config

import "time"

// EngineType represents the type of translation engine
type EngineType string

const (
	EngineInternal      EngineType = "internal"
	EngineTerminalAgent EngineType = "terminal-agent"
	EngineOllama        EngineType = "ollama"
)

// TerminalAgentType represents the terminal agent type
type TerminalAgentType string

const (
	AgentClaudeCode TerminalAgentType = "claude-code"
	AgentGeminiCLI  TerminalAgentType = "gemini-cli"
	AgentCodex      TerminalAgentType = "codex"
)

// EngineConfig holds translation engine settings
type EngineConfig struct {
	Type          EngineType          `json:"type"`
	Internal      InternalConfig      `json:"internal"`
	TerminalAgent TerminalAgentConfig `json:"terminalAgent"`
	Ollama        OllamaConfig        `json:"ollama"`
}

// InternalConfig holds internal (Yzma) engine settings
type InternalConfig struct {
	ModelPath   string `json:"modelPath"`
	ContextSize int    `json:"contextSize"`
}

// TerminalAgentConfig holds terminal agent settings
type TerminalAgentConfig struct {
	Selected   TerminalAgentType   `json:"selected"`
	ClaudeCode TerminalAgentOption `json:"claudeCode"`
	GeminiCLI  TerminalAgentOption `json:"geminiCli"`
	Codex      TerminalAgentOption `json:"codex"`
}

// TerminalAgentOption holds settings for a terminal agent
type TerminalAgentOption struct {
	Executable string   `json:"executable"` // path to executable (empty = use PATH)
	Args       []string `json:"args"`       // additional arguments
	Timeout    int      `json:"timeout"`    // seconds
}

// OllamaConfig holds Ollama engine settings
type OllamaConfig struct {
	Host    string `json:"host"`
	Model   string `json:"model"`
	Timeout int    `json:"timeout"` // seconds
}

// DefaultEngineConfig returns default engine settings
func DefaultEngineConfig() EngineConfig {
	return EngineConfig{
		Type: EngineInternal,
		Internal: InternalConfig{
			ContextSize: 2048,
		},
		TerminalAgent: TerminalAgentConfig{
			Selected: AgentClaudeCode,
			ClaudeCode: TerminalAgentOption{
				Executable: "claude",
				Timeout:    60,
			},
			GeminiCLI: TerminalAgentOption{
				Executable: "gemini",
				Timeout:    60,
			},
			Codex: TerminalAgentOption{
				Executable: "codex",
				Timeout:    60,
			},
		},
		Ollama: OllamaConfig{
			Host:    "http://localhost:11434",
			Model:   "llama3.2",
			Timeout: 120,
		},
	}
}

// SetEngineType sets the engine type with validation
func (c *Config) SetEngineType(engine EngineType) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch engine {
	case EngineInternal, EngineTerminalAgent, EngineOllama:
		c.Engine.Type = engine
	default:
		c.Engine.Type = EngineInternal
	}
}

// SetTerminalAgent sets the selected terminal agent with validation
func (c *Config) SetTerminalAgent(agent TerminalAgentType) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch agent {
	case AgentClaudeCode, AgentGeminiCLI, AgentCodex:
		c.Engine.TerminalAgent.Selected = agent
	default:
		c.Engine.TerminalAgent.Selected = AgentClaudeCode
	}
}

// GetSelectedTerminalAgent returns the config for the selected terminal agent
func (c *Config) GetSelectedTerminalAgent() TerminalAgentOption {
	c.mu.RLock()
	defer c.mu.RUnlock()

	switch c.Engine.TerminalAgent.Selected {
	case AgentClaudeCode:
		return c.Engine.TerminalAgent.ClaudeCode
	case AgentGeminiCLI:
		return c.Engine.TerminalAgent.GeminiCLI
	case AgentCodex:
		return c.Engine.TerminalAgent.Codex
	default:
		return c.Engine.TerminalAgent.ClaudeCode
	}
}

// GetSelectedTerminalAgentTimeout returns the timeout for the selected terminal agent
func (c *Config) GetSelectedTerminalAgentTimeout() time.Duration {
	agent := c.GetSelectedTerminalAgent()
	return time.Duration(agent.Timeout) * time.Second
}

// GetSelectedTerminalAgentExecutable returns the executable path for the selected terminal agent
func (c *Config) GetSelectedTerminalAgentExecutable() string {
	agent := c.GetSelectedTerminalAgent()
	return agent.Executable
}

// GetSelectedTerminalAgentArgs returns the additional arguments for the selected terminal agent
func (c *Config) GetSelectedTerminalAgentArgs() []string {
	agent := c.GetSelectedTerminalAgent()
	if agent.Args == nil {
		return nil
	}
	args := make([]string, len(agent.Args))
	copy(args, agent.Args)
	return args
}

// SetOllamaHost sets the Ollama host URL
func (c *Config) SetOllamaHost(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Engine.Ollama.Host = host
}

// SetOllamaModel sets the Ollama model name
func (c *Config) SetOllamaModel(model string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Engine.Ollama.Model = model
}

// SetInternalModelPath sets the internal model path
func (c *Config) SetInternalModelPath(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Engine.Internal.ModelPath = path
}
