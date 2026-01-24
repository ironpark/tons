package engine

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
		Args:    []string{"--model", "haiku", "--tools", "", "--output-format", "stream-json", "--verbose", "--include-partial-messages", "-p"},
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

// buildArgs constructs command arguments with optional system prompt support
func (e *TerminalEngine) buildArgs(prompt, systemPrompt string) []string {
	args := make([]string, len(e.config.Args))
	copy(args, e.config.Args)

	// For Claude Code, add system prompt before -p flag if provided
	if e.name == string(TerminalClaudeCode) && systemPrompt != "" {
		args = append(args, "--system-prompt", systemPrompt)
	}

	args = append(args, prompt)
	return args
}

// Translate performs non-streaming translation
func (e *TerminalEngine) Translate(ctx context.Context, req Request) (Response, error) {
	if req.Text == "" {
		return Response{Text: "", Done: true}, nil
	}

	prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

	ctx, cancel := context.WithTimeout(ctx, e.config.Timeout)
	defer cancel()

	args := e.buildArgs(prompt, req.SystemPrompt)
	cmd := exec.CommandContext(ctx, e.config.Command, args...)
	slog.Info("Translate", "cmd", cmd)
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return Response{}, fmt.Errorf("translation timed out")
		}
		return Response{}, fmt.Errorf("terminal agent error: %w", err)
	}

	return Response{Text: strings.TrimSpace(string(output)), Done: true}, nil
}

// claudeCodeEvent represents the JSON structure from Claude Code stream output
type claudeCodeEvent struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype,omitempty"`
	Event   *struct {
		Type  string `json:"type"`
		Delta *struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"delta,omitempty"`
	} `json:"event,omitempty"`
	Result string `json:"result,omitempty"`
}

// TranslateStream performs streaming translation
func (e *TerminalEngine) TranslateStream(ctx context.Context, req Request) (<-chan Response, error) {
	ch := make(chan Response)

	go func() {
		defer close(ch)

		if req.Text == "" {
			ch <- Response{Text: "", Done: true}
			return
		}

		prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)
		slog.Info("TranslateStream", "prompt", prompt)

		ctx, cancel := context.WithTimeout(ctx, e.config.Timeout)
		defer cancel()

		args := e.buildArgs(prompt, req.SystemPrompt)
		cmd := exec.CommandContext(ctx, e.config.Command, args...)
		slog.Info("TranslateStream", "cmd", cmd)

		stdout, err := cmd.StdoutPipe()
		cmd.Stderr = os.Stderr
		if err != nil {
			ch <- ErrorResponsef("failed to create pipe: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			ch <- ErrorResponsef("failed to start command: %v", err)
			return
		}

		// Check if this is a Claude Code engine that outputs JSON stream
		isClaudeCode := e.name == string(TerminalClaudeCode)

		if isClaudeCode {
			e.streamClaudeCodeOutput(ctx, cmd, stdout, ch)
		} else {
			e.streamRawOutput(ctx, cmd, stdout, ch)
		}
		cmd.Wait()

	}()

	return ch, nil
}

// streamClaudeCodeOutput handles JSON streaming output from Claude Code CLI
func (e *TerminalEngine) streamClaudeCodeOutput(ctx context.Context, cmd *exec.Cmd, stdout io.ReadCloser, ch chan<- Response) {
	// lineResult holds the result of reading a line
	type lineResult struct {
		line string
		err  error
	}

	lineCh := make(chan lineResult)
	go func() {
		defer close(lineCh)
		scanner := bufio.NewScanner(stdout)
		// // Increase buffer size for potentially large JSON lines
		// scanner.Buffer(make([]byte, 64*1024), 1024*1024)
		for scanner.Scan() {
			text := scanner.Text()
			slog.Info("streamClaudeCodeOutput", "LINE", text)
			lineCh <- lineResult{line: text}
		}
		if err := scanner.Err(); err != nil {
			lineCh <- lineResult{err: err}
		}
	}()
	text := ""
	for {
		select {
		case <-ctx.Done():
			gracefulShutdown(cmd.Process)
			ch <- ErrorResponse("translation timed out")
			return
		case result := <-lineCh:
			// slog.Info("DD", "LINE", result.line)
			// if !ok {
			// 	cmd.Wait()
			// 	ch <- Response{Done: true}
			// 	return
			// }
			if result.err != nil {
				ch <- ErrorResponsef("read error: %v", result.err)
				cmd.Wait()
				return
			}

			// Parse JSON line
			var event claudeCodeEvent
			if err := json.Unmarshal([]byte(result.line), &event); err != nil {
				// Skip non-JSON lines
				continue
			}

			switch event.Type {
			case "stream_event":
				// Extract text from content_block_delta events
				if event.Event != nil && event.Event.Type == "content_block_delta" && event.Event.Delta != nil {
					if event.Event.Delta.Type == "text_delta" && event.Event.Delta.Text != "" {
						text += event.Event.Delta.Text
						ch <- Response{Text: text, Done: false}
					}
				}
			case "result":
				// Final result - we already streamed the content, just mark as done
				cmd.Wait()
				ch <- Response{Done: true}
				return
			}
		}
	}
}

// streamRawOutput handles raw byte streaming for non-Claude Code engines
func (e *TerminalEngine) streamRawOutput(ctx context.Context, cmd *exec.Cmd, stdout io.ReadCloser, ch chan<- Response) {
	// readResult holds the result of a read operation
	type readResult struct {
		data []byte
		err  error
	}

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

	for {
		select {
		case <-ctx.Done():
			gracefulShutdown(cmd.Process)
			ch <- ErrorResponse("translation timed out")
			return
		case result, ok := <-readCh:
			if !ok {
				cmd.Wait()
				ch <- Response{Done: true}
				return
			}
			if result.err != nil {
				ch <- ErrorResponsef("read error: %v", result.err)
				cmd.Wait()
				return
			}
			ch <- Response{Text: string(result.data), Done: false}
		}
	}
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
