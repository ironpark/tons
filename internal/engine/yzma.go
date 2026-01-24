package engine

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/hybridgroup/yzma/pkg/llama"
)

const (
	defaultNCtx = 2048
)

// Yzma is the local LLM translation engine using yzma/llama
type Yzma struct {
	ModelPath   string
	Sampling    SamplingConfig
	ContextSize int
	model       llama.Model
	vocab       llama.Vocab
	mu          sync.Mutex
	initialized bool
	inUse       chan struct{} // semaphore for inference concurrency control
}

// YzmaOption is a functional option for configuring Yzma
type YzmaOption func(*Yzma)

// WithYzmaSampling sets the sampling configuration
func WithYzmaSampling(cfg SamplingConfig) YzmaOption {
	return func(y *Yzma) {
		y.Sampling = cfg
	}
}

// WithYzmaContextSize sets the context size for inference
func WithYzmaContextSize(size int) YzmaOption {
	return func(y *Yzma) {
		y.ContextSize = size
	}
}

// NewYzma creates a new Yzma engine with the given model path and options
func NewYzma(modelPath string, opts ...YzmaOption) *Yzma {
	y := &Yzma{
		ModelPath:   modelPath,
		Sampling:    DefaultSamplingConfig(),
		ContextSize: defaultNCtx,
		inUse:       make(chan struct{}, 1),
	}
	for _, opt := range opts {
		opt(y)
	}
	return y
}

// Name returns the engine name
func (e *Yzma) Name() string {
	return "yzma"
}

// Available checks if model file exists
func (e *Yzma) Available() bool {
	_, err := os.Stat(e.ModelPath)
	return err == nil
}

// Initialize loads the model (lazy initialization)
func (e *Yzma) Initialize() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.initialized {
		return nil
	}

	llama.Init()

	params := llama.ModelDefaultParams()
	model, err := llama.ModelLoadFromFile(e.ModelPath, params)
	if err != nil {
		return err
	}

	e.model = model
	e.vocab = llama.ModelGetVocab(model)
	e.initialized = true
	return nil
}

// Close releases model resources
func (e *Yzma) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.model != 0 {
		llama.ModelFree(e.model)
		e.model = 0
		e.vocab = 0
		e.initialized = false
	}
	return nil
}

// acquireModel acquires exclusive access to the model for inference.
// Returns a release function that must be called when done.
func (e *Yzma) acquireModel(ctx context.Context) (release func(), err error) {
	select {
	case e.inUse <- struct{}{}:
		return func() { <-e.inUse }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// generationCallback is called for each generated token piece
type generationCallback func(piece string) bool

// generateTokens handles the common token generation logic
func (e *Yzma) generateTokens(ctx context.Context, prompt string, cb generationCallback) error {
	// Create context for inference
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = uint32(e.ContextSize)
	llamaCtx, err := llama.InitFromModel(e.model, ctxParams)
	if err != nil {
		return err
	}
	defer llama.Free(llamaCtx)

	// Tokenize the prompt
	tokens := llama.Tokenize(e.vocab, prompt, true, false)

	// Create sampler chain using config
	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(sampler, llama.SamplerInitTempExt(e.Sampling.Temperature, 0, 1))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopP(e.Sampling.TopP, 1))
	llama.SamplerChainAdd(sampler, llama.SamplerInitDist(0))
	defer llama.SamplerFree(sampler)

	// Process initial prompt
	batch := llama.BatchGetOne(tokens)
	if _, err := llama.Decode(llamaCtx, batch); err != nil {
		return err
	}

	// Generate response tokens
	eosToken := llama.VocabEOS(e.vocab)
	buf := make([]byte, 256)

	for range e.Sampling.MaxTokens {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Sample next token
		token := llama.SamplerSample(sampler, llamaCtx, -1)

		// Check for end of sequence
		if token == eosToken {
			return nil
		}

		// Convert token to text
		n := llama.TokenToPiece(e.vocab, token, buf, 0, false)
		if n > 0 {
			piece := string(buf[:n])
			if !cb(piece) {
				return nil // Callback requested stop
			}
		}

		// Prepare next batch with the new token
		batch = llama.BatchGetOne([]llama.Token{token})
		if _, err := llama.Decode(llamaCtx, batch); err != nil {
			return err
		}
	}

	return nil // Max tokens reached
}

// Translate performs translation (non-streaming)
func (e *Yzma) Translate(ctx context.Context, req Request) (Response, error) {
	if req.Text == "" {
		return Response{Text: "", Done: true}, nil
	}

	if err := e.Initialize(); err != nil {
		return Response{}, fmt.Errorf("yzma error: %w", err)
	}

	prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

	release, err := e.acquireModel(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("yzma error: failed to acquire model: %w", err)
	}
	defer release()

	var result strings.Builder

	err = e.generateTokens(ctx, prompt, func(piece string) bool {
		result.WriteString(piece)
		return true
	})

	if err != nil {
		// If we have partial results, return them along with the error
		if result.Len() > 0 {
			return Response{Text: result.String(), Done: true}, fmt.Errorf("yzma error: %w", err)
		}
		return Response{}, fmt.Errorf("yzma error: %w", err)
	}

	return Response{Text: result.String(), Done: true}, nil
}

// TranslateStream performs streaming translation
func (e *Yzma) TranslateStream(ctx context.Context, req Request) (<-chan Response, error) {
	if err := e.Initialize(); err != nil {
		return nil, fmt.Errorf("yzma error: %w", err)
	}

	ch := make(chan Response)

	go func() {
		defer close(ch)

		if req.Text == "" {
			ch <- Response{Text: "", Done: true}
			return
		}

		prompt := BuildPrompt(req.Prompt, req.Text, req.SourceLang, req.TargetLang)

		release, err := e.acquireModel(ctx)
		if err != nil {
			ch <- ErrorResponsef("yzma error: failed to acquire model: %v", err)
			return
		}
		defer release()

		err = e.generateTokens(ctx, prompt, func(piece string) bool {
			select {
			case ch <- Response{Text: piece, Done: false}:
				return true
			case <-ctx.Done():
				return false
			}
		})

		if err != nil {
			ch <- ErrorResponsef("yzma error: %v", err)
			return
		}

		ch <- Response{Done: true}
	}()

	return ch, nil
}
