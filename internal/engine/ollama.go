package engine

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

// Ollama uses Ollama for translation
type Ollama struct {
	Host     string
	Model    string
	Timeout  time.Duration
	Sampling SamplingConfig
	client   *api.Client
}

// OllamaOption is a functional option for configuring Ollama
type OllamaOption func(*Ollama)

// WithOllamaHost sets the Ollama host URL
func WithOllamaHost(host string) OllamaOption {
	return func(o *Ollama) {
		o.Host = host
	}
}

// WithOllamaTimeout sets the request timeout
func WithOllamaTimeout(timeout time.Duration) OllamaOption {
	return func(o *Ollama) {
		o.Timeout = timeout
	}
}

// WithOllamaSampling sets the sampling configuration
func WithOllamaSampling(cfg SamplingConfig) OllamaOption {
	return func(o *Ollama) {
		o.Sampling = cfg
	}
}

// NewOllama creates a new Ollama engine with optional configuration
func NewOllama(model string, opts ...OllamaOption) *Ollama {
	o := &Ollama{
		Host:     "http://localhost:11434",
		Model:    model,
		Timeout:  120 * time.Second,
		Sampling: DefaultSamplingConfig(),
	}
	for _, opt := range opts {
		opt(o)
	}

	// Create API client
	hostURL, err := url.Parse(o.Host)
	if err != nil {
		hostURL, _ = url.Parse("http://localhost:11434")
	}
	o.client = api.NewClient(hostURL, &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
			MaxIdleConnsPerHost: 5,
		},
	})

	return o
}

// NewOllamaWithHost creates a new Ollama engine with custom host
// Deprecated: Use NewOllama with WithOllamaHost option instead
func NewOllamaWithHost(host, model string) *Ollama {
	return NewOllama(model, WithOllamaHost(host))
}

// Name returns the engine name
func (e *Ollama) Name() string {
	return "ollama:" + e.Model
}

// Available checks if Ollama is accessible
func (e *Ollama) Available() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := e.client.List(ctx)
	return err == nil
}

// Close releases resources held by the Ollama engine
func (e *Ollama) Close() error {
	return nil
}

// Translate performs translation using Ollama (non-streaming)
func (e *Ollama) Translate(ctx context.Context, req Request) Response {
	if req.Text == "" {
		return Response{Text: "", Done: true}
	}

	prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

	ctx, cancel := context.WithTimeout(ctx, e.Timeout)
	defer cancel()

	genReq := &api.GenerateRequest{
		Model:  e.Model,
		Prompt: prompt,
	}

	var result strings.Builder
	err := e.client.Generate(ctx, genReq, func(resp api.GenerateResponse) error {
		result.WriteString(resp.Response)
		return nil
	})

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return ErrorResponse("Translation timed out")
		}
		return ErrorResponsef("Ollama error: %v", err)
	}

	return Response{Text: strings.TrimSpace(result.String()), Done: true}
}

// TranslateStream performs streaming translation using Ollama
func (e *Ollama) TranslateStream(ctx context.Context, req Request) <-chan Response {
	ch := make(chan Response)

	go func() {
		defer close(ch)

		if req.Text == "" {
			ch <- Response{Text: "", Done: true}
			return
		}

		prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

		ctx, cancel := context.WithTimeout(ctx, e.Timeout)
		defer cancel()

		genReq := &api.GenerateRequest{
			Model:  e.Model,
			Prompt: prompt,
		}

		err := e.client.Generate(ctx, genReq, func(resp api.GenerateResponse) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				ch <- Response{Text: resp.Response, Done: resp.Done}
				return nil
			}
		})

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				ch <- ErrorResponse("Translation timed out")
			} else {
				ch <- ErrorResponsef("Ollama error: %v", err)
			}
		}
	}()

	return ch
}

// Model represents an Ollama model
type Model struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
}

// GetModels returns available Ollama models
func (e *Ollama) GetModels() ([]Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listResp, err := e.client.List(ctx)
	if err != nil {
		return nil, err
	}

	models := make([]Model, len(listResp.Models))
	for i, m := range listResp.Models {
		models[i] = Model{
			Name:       m.Name,
			ModifiedAt: m.ModifiedAt.Format(time.RFC3339),
			Size:       m.Size,
		}
	}

	return models, nil
}

// OllamaModels returns available Ollama model names from a host
func OllamaModels(host string) []string {
	eng := NewOllama("", WithOllamaHost(host))
	models, err := eng.GetModels()
	if err != nil {
		return []string{}
	}

	names := make([]string, len(models))
	for i, m := range models {
		names[i] = m.Name
	}
	return names
}
