package engine

import (
	"context"
	"time"
)

// ClaudeCode uses Claude Code CLI for translation
type ClaudeCode struct {
	config TerminalConfig
}

// ClaudeCodeOption is a functional option for ClaudeCode
type ClaudeCodeOption func(*ClaudeCode)

// WithClaudeCodeTimeout sets the timeout for ClaudeCode operations
func WithClaudeCodeTimeout(timeout time.Duration) ClaudeCodeOption {
	return func(c *ClaudeCode) {
		c.config.Timeout = timeout
	}
}

// NewClaudeCode creates a new Claude Code engine
func NewClaudeCode(opts ...ClaudeCodeOption) *ClaudeCode {
	c := &ClaudeCode{
		config: TerminalConfig{
			Command: "claude",
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
func (e *ClaudeCode) Name() string {
	return "claude-code"
}

// Available checks if Claude Code CLI is installed
func (e *ClaudeCode) Available() bool {
	return terminalAvailable(e.config.Command)
}

// Close releases resources held by the ClaudeCode engine
func (e *ClaudeCode) Close() error {
	return nil
}

// Translate performs translation using Claude Code CLI
func (e *ClaudeCode) Translate(ctx context.Context, req Request) Response {
	return terminalTranslate(ctx, e.config, req)
}

// TranslateStream performs streaming translation using Claude Code CLI
func (e *ClaudeCode) TranslateStream(ctx context.Context, req Request) <-chan Response {
	return terminalTranslateStream(ctx, e.config, req)
}
