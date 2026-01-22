package engine

import (
	"context"
	"fmt"
	"strings"
)

// Request represents a translation request
type Request struct {
	Text       string `json:"text"`
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
	Prompt     string `json:"prompt"`
}

// Response represents a translation response
type Response struct {
	Text  string `json:"text"`
	Done  bool   `json:"done"`
	Error string `json:"error,omitempty"`
}

// ErrorResponse creates an error response with the given message
func ErrorResponse(err string) Response {
	return Response{Error: err, Done: true}
}

// ErrorResponsef creates an error response with a formatted message
func ErrorResponsef(format string, args ...any) Response {
	return Response{Error: fmt.Sprintf(format, args...), Done: true}
}

// SamplingConfig holds sampling parameters for LLM generation
type SamplingConfig struct {
	Temperature float32
	TopP        float32
	MaxTokens   int
}

// DefaultSamplingConfig returns default sampling parameters
func DefaultSamplingConfig() SamplingConfig {
	return SamplingConfig{
		Temperature: 0.7,
		TopP:        0.9,
		MaxTokens:   512,
	}
}

// Engine is the interface for translation engines
type Engine interface {
	Name() string
	Translate(ctx context.Context, req Request) Response
	TranslateStream(ctx context.Context, req Request) <-chan Response
	Available() bool
	Close() error
}

// BuildPrompt replaces template variables with actual values
func BuildPrompt(template, text, sourceLang, targetLang string) string {
	prompt := template
	prompt = strings.ReplaceAll(prompt, "{{text}}", text)
	prompt = strings.ReplaceAll(prompt, "{{source_lang}}", sourceLang)
	prompt = strings.ReplaceAll(prompt, "{{target_lang}}", targetLang)
	return prompt
}
