package engine

import (
	"context"
	"time"
)

// Codex uses Codex CLI for translation
type Codex struct {
	config TerminalConfig
}

// CodexOption is a functional option for Codex
type CodexOption func(*Codex)

// WithCodexTimeout sets the timeout for Codex operations
func WithCodexTimeout(timeout time.Duration) CodexOption {
	return func(c *Codex) {
		c.config.Timeout = timeout
	}
}

// NewCodex creates a new Codex engine
func NewCodex(opts ...CodexOption) *Codex {
	c := &Codex{
		config: TerminalConfig{
			Command: "codex",
			Args:    []string{"-p"},
			Timeout: 60 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Name returns the engine name
func (e *Codex) Name() string {
	return "codex"
}

// Available checks if Codex CLI is installed
func (e *Codex) Available() bool {
	return terminalAvailable(e.config.Command)
}

// Close releases resources held by the Codex engine
func (e *Codex) Close() error {
	return nil
}

// Translate performs translation using Codex CLI
func (e *Codex) Translate(ctx context.Context, req Request) Response {
	return terminalTranslate(ctx, e.config, req)
}

// TranslateStream performs streaming translation using Codex CLI
func (e *Codex) TranslateStream(ctx context.Context, req Request) <-chan Response {
	return terminalTranslateStream(ctx, e.config, req)
}
