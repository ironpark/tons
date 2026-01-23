package config

// Theme represents the application theme
type Theme string

const (
	ThemeLight  Theme = "light"
	ThemeDark   Theme = "dark"
	ThemeSystem Theme = "system"
)

// GeneralConfig holds general application settings
type GeneralConfig struct {
	Theme    Theme  `json:"theme"`
	Language string `json:"language"`
}

// DefaultGeneralConfig returns default general settings
func DefaultGeneralConfig() GeneralConfig {
	return GeneralConfig{
		Theme:    ThemeSystem,
		Language: "system",
	}
}

// SetTheme sets the theme with validation
func (c *Config) SetTheme(theme Theme) {
	switch theme {
	case ThemeLight, ThemeDark, ThemeSystem:
		c.General.Theme = theme
	default:
		c.General.Theme = ThemeSystem
	}
}

// SetLanguage sets the language
func (c *Config) SetLanguage(lang string) {
	c.General.Language = lang
}
