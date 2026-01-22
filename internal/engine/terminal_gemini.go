package engine

import (
	"context"
	"time"
)

// GeminiCLI uses Gemini CLI for translation
type GeminiCLI struct {
	config TerminalConfig
}

// GeminiCLIOption is a functional option for GeminiCLI
type GeminiCLIOption func(*GeminiCLI)

// WithGeminiCLITimeout sets the timeout for GeminiCLI operations
func WithGeminiCLITimeout(timeout time.Duration) GeminiCLIOption {
	return func(g *GeminiCLI) {
		g.config.Timeout = timeout
	}
}

// NewGeminiCLI creates a new Gemini CLI engine
func NewGeminiCLI(opts ...GeminiCLIOption) *GeminiCLI {
	g := &GeminiCLI{
		config: TerminalConfig{
			Command: "gemini",
			Args:    []string{"-p"},
			Timeout: 60 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Name returns the engine name
func (e *GeminiCLI) Name() string {
	return "gemini-cli"
}

// Available checks if Gemini CLI is installed
func (e *GeminiCLI) Available() bool {
	return terminalAvailable(e.config.Command)
}

// Close releases resources held by the GeminiCLI engine
func (e *GeminiCLI) Close() error {
	return nil
}

// Translate performs translation using Gemini CLI
func (e *GeminiCLI) Translate(ctx context.Context, req Request) Response {
	return terminalTranslate(ctx, e.config, req)
}

// TranslateStream performs streaming translation using Gemini CLI
func (e *GeminiCLI) TranslateStream(ctx context.Context, req Request) <-chan Response {
	return terminalTranslateStream(ctx, e.config, req)
}
