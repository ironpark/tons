package config

// DefaultPrompt is the default translation prompt template
const DefaultPrompt = `Translate the following text from {{source_lang}} to {{target_lang}}.
Keep the original formatting and tone.
Only return the translated text without any explanations.

Text to translate:
{{text}}`

// PromptConfig holds prompt settings
type PromptConfig struct {
	Template string `json:"template"`
}

// DefaultPromptConfig returns default prompt settings
func DefaultPromptConfig() PromptConfig {
	return PromptConfig{
		Template: DefaultPrompt,
	}
}

// SetPrompt sets the prompt template
func (c *Config) SetPrompt(prompt string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Prompt.Template = prompt
}

// ResetPrompt restores the prompt to default
func (c *Config) ResetPrompt() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Prompt.Template = DefaultPrompt
}
