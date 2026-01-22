package engine

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"time"
)

// TerminalConfig holds common configuration for terminal-based engines
type TerminalConfig struct {
	Command string        // CLI command name (e.g., "claude", "gemini")
	Args    []string      // Base arguments before prompt
	Timeout time.Duration
}

// readResult holds the result of a read operation
type readResult struct {
	data []byte
	err  error
}

// terminalTranslate performs non-streaming translation using a CLI command
func terminalTranslate(ctx context.Context, cfg TerminalConfig, req Request) Response {
	if req.Text == "" {
		return Response{Text: "", Done: true}
	}

	prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	args := append(cfg.Args, prompt)
	cmd := exec.CommandContext(ctx, cfg.Command, args...)
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return ErrorResponse("Translation timed out")
		}
		return ErrorResponsef("Terminal agent error: %v", err)
	}

	return Response{Text: strings.TrimSpace(string(output)), Done: true}
}

// terminalTranslateStream performs streaming translation using a CLI command
func terminalTranslateStream(ctx context.Context, cfg TerminalConfig, req Request) <-chan Response {
	ch := make(chan Response)

	go func() {
		defer close(ch)

		if req.Text == "" {
			ch <- Response{Text: "", Done: true}
			return
		}

		prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()

		args := append(cfg.Args, prompt)
		cmd := exec.CommandContext(ctx, cfg.Command, args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ch <- ErrorResponsef("Failed to create pipe: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			ch <- ErrorResponsef("Failed to start command: %v", err)
			return
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

		// Main loop properly handles context cancellation
		for {
			select {
			case <-ctx.Done():
				cmd.Process.Kill()
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

// terminalAvailable checks if a command is available in PATH
func terminalAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// AvailableTerminalEngines returns all installed terminal-based engines
func AvailableTerminalEngines() []Engine {
	var available []Engine
	if e := NewClaudeCode(); e.Available() {
		available = append(available, e)
	}
	if e := NewGeminiCLI(); e.Available() {
		available = append(available, e)
	}
	if e := NewCodex(); e.Available() {
		available = append(available, e)
	}
	return available
}
