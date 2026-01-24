package config

// DefaultPrompt is the default translation prompt template
const DefaultPrompt = `Translate the following text from {{source_lang}} to {{target_lang}}.
Keep the original formatting and tone.
Only return the translated text without any explanations.

Text to translate:
{{text}}`

// DefaultSystemPrompt is the default system prompt for translation
const DefaultSystemPrompt = `You are a professional translator. Translate accurately while preserving the original tone, style, and formatting. Only output the translation without explanations.`

// PromptConfig holds prompt settings
type PromptConfig struct {
	Template     string `json:"template"`
	SystemPrompt string `json:"systemPrompt"`
}

// DefaultPromptConfig returns default prompt settings
func DefaultPromptConfig() PromptConfig {
	return PromptConfig{
		Template:     DefaultPrompt,
		SystemPrompt: DefaultSystemPrompt,
	}
}

// SetPrompt sets the prompt template
func (c *Config) SetPrompt(prompt string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Prompt.Template = prompt
}

// SetSystemPrompt sets the system prompt
func (c *Config) SetSystemPrompt(systemPrompt string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Prompt.SystemPrompt = systemPrompt
}

// ResetPrompt restores the prompt to default
func (c *Config) ResetPrompt() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Prompt.Template = DefaultPrompt
	c.Prompt.SystemPrompt = DefaultSystemPrompt
}

// SetPromptConfig sets the entire prompt config
func (c *Config) SetPromptConfig(prompt PromptConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Prompt = prompt
}
