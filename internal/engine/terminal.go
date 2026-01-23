package engine

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// TerminalEngineType represents predefined terminal engine types
type TerminalEngineType string

const (
	TerminalClaudeCode TerminalEngineType = "claude-code"
	TerminalGeminiCLI  TerminalEngineType = "gemini-cli"
	TerminalCodex      TerminalEngineType = "codex"
)

// TerminalConfig holds configuration for terminal-based engines
type TerminalConfig struct {
	Command string        // CLI command name (e.g., "claude", "gemini")
	Args    []string      // Base arguments before prompt
	Timeout time.Duration // Timeout for translation operations
}

// predefinedEngines contains default configurations for known terminal engines
var predefinedEngines = map[TerminalEngineType]TerminalConfig{
	TerminalClaudeCode: {
		Command: "claude",
		Args:    []string{"-p"},
		Timeout: 60 * time.Second,
	},
	TerminalGeminiCLI: {
		Command: "gemini",
		Args:    []string{"-p"},
		Timeout: 60 * time.Second,
	},
	TerminalCodex: {
		Command: "codex",
		Args:    []string{"-p"},
		Timeout: 60 * time.Second,
	},
}

// TerminalEngine is a unified engine for CLI-based translation tools
type TerminalEngine struct {
	name   string
	config TerminalConfig
}

// TerminalEngineOption is a functional option for TerminalEngine
type TerminalEngineOption func(*TerminalEngine)

// WithTerminalTimeout sets the timeout for terminal operations
func WithTerminalTimeout(timeout time.Duration) TerminalEngineOption {
	return func(e *TerminalEngine) {
		e.config.Timeout = timeout
	}
}

// WithTerminalArgs sets additional arguments for the terminal command
func WithTerminalArgs(args []string) TerminalEngineOption {
	return func(e *TerminalEngine) {
		e.config.Args = args
	}
}

// WithTerminalCommand overrides the command to execute
func WithTerminalCommand(command string) TerminalEngineOption {
	return func(e *TerminalEngine) {
		e.config.Command = command
	}
}

// NewTerminalEngine creates a new terminal engine from a predefined type
func NewTerminalEngine(engineType TerminalEngineType, opts ...TerminalEngineOption) *TerminalEngine {
	cfg, ok := predefinedEngines[engineType]
	if !ok {
		// Default fallback configuration
		cfg = TerminalConfig{
			Command: string(engineType),
			Args:    []string{"-p"},
			Timeout: 60 * time.Second,
		}
	}

	e := &TerminalEngine{
		name:   string(engineType),
		config: cfg,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// NewCustomTerminalEngine creates a terminal engine with custom configuration
func NewCustomTerminalEngine(name, command string, args []string, opts ...TerminalEngineOption) *TerminalEngine {
	e := &TerminalEngine{
		name: name,
		config: TerminalConfig{
			Command: command,
			Args:    args,
			Timeout: 60 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Name returns the engine name
func (e *TerminalEngine) Name() string {
	return e.name
}

// Available checks if the CLI command is available in PATH
func (e *TerminalEngine) Available() bool {
	_, err := exec.LookPath(e.config.Command)
	return err == nil
}

// Close releases resources (no-op for terminal engines)
func (e *TerminalEngine) Close() error {
	return nil
}

// Translate performs non-streaming translation
func (e *TerminalEngine) Translate(ctx context.Context, req Request) Response {
	if req.Text == "" {
		return Response{Text: "", Done: true}
	}

	prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

	ctx, cancel := context.WithTimeout(ctx, e.config.Timeout)
	defer cancel()

	args := append(e.config.Args, prompt)
	cmd := exec.CommandContext(ctx, e.config.Command, args...)
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return ErrorResponse("Translation timed out")
		}
		return ErrorResponsef("Terminal agent error: %v", err)
	}

	return Response{Text: strings.TrimSpace(string(output)), Done: true}
}

// TranslateStream performs streaming translation
func (e *TerminalEngine) TranslateStream(ctx context.Context, req Request) <-chan Response {
	ch := make(chan Response)

	go func() {
		defer close(ch)

		if req.Text == "" {
			ch <- Response{Text: "", Done: true}
			return
		}

		prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

		ctx, cancel := context.WithTimeout(ctx, e.config.Timeout)
		defer cancel()

		args := append(e.config.Args, prompt)
		cmd := exec.CommandContext(ctx, e.config.Command, args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ch <- ErrorResponsef("Failed to create pipe: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			ch <- ErrorResponsef("Failed to start command: %v", err)
			return
		}

		// readResult holds the result of a read operation
		type readResult struct {
			data []byte
			err  error
		}

		// Use a separate goroutine for reading to avoid blocking on stdout.Read()
		readCh := make(chan readResult)
		go func() {
			defer close(readCh)
			buf := make([]byte, 1024)
			for {
				n, err := stdout.Read(buf)
				if n > 0 {
					data := make([]byte, n)
					copy(data, buf[:n])
					readCh <- readResult{data: data}
				}
				if err != nil {
					if err != io.EOF {
						readCh <- readResult{err: err}
					}
					return
				}
			}
		}()

		// Main loop with proper context cancellation and graceful shutdown
		for {
			select {
			case <-ctx.Done():
				gracefulShutdown(cmd.Process)
				ch <- ErrorResponse("Translation timed out")
				return
			case result, ok := <-readCh:
				if !ok {
					// Read goroutine finished, wait for command and send done
					cmd.Wait()
					ch <- Response{Done: true}
					return
				}
				if result.err != nil {
					ch <- ErrorResponsef("Read error: %v", result.err)
					cmd.Wait()
					return
				}
				ch <- Response{Text: string(result.data), Done: false}
			}
		}
	}()

	return ch
}

// gracefulShutdown attempts to terminate a process gracefully before force killing
func gracefulShutdown(proc *os.Process) {
	if proc == nil {
		return
	}

	// First, try SIGTERM for graceful shutdown
	proc.Signal(syscall.SIGTERM)

	// Wait up to 3 seconds for graceful termination
	done := make(chan struct{})
	go func() {
		proc.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Process terminated gracefully
		return
	case <-time.After(3 * time.Second):
		// Force kill if still running
		proc.Kill()
		proc.Wait()
	}
}

// AvailableTerminalEngines returns all installed terminal-based engines
func AvailableTerminalEngines() []Engine {
	var available []Engine

	for engineType := range predefinedEngines {
		e := NewTerminalEngine(engineType)
		if e.Available() {
			available = append(available, e)
		}
	}

	return available
}

// Convenience constructors for backward compatibility
// These can be removed once all callers are updated

// NewClaudeCode creates a Claude Code terminal engine
func NewClaudeCode(opts ...TerminalEngineOption) *TerminalEngine {
	return NewTerminalEngine(TerminalClaudeCode, opts...)
}

// NewGeminiCLI creates a Gemini CLI terminal engine
func NewGeminiCLI(opts ...TerminalEngineOption) *TerminalEngine {
	return NewTerminalEngine(TerminalGeminiCLI, opts...)
}

// NewCodex creates a Codex terminal engine
func NewCodex(opts ...TerminalEngineOption) *TerminalEngine {
	return NewTerminalEngine(TerminalCodex, opts...)
}
