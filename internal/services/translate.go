package services

import (
	"context"
	"log/slog"

	"github.com/ironpark/tons/internal/config"
	"github.com/ironpark/tons/internal/engine"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type TranslateService struct {
	cfg *config.Config
	app *application.App
}

func NewTranslateService(cfg *config.Config) *TranslateService {
	return &TranslateService{
		cfg: cfg,
	}
}

func (ts *TranslateService) Translate(sourceLang, targetLang, text string) error {
	snapshot := ts.cfg.Snapshot()
	engineCfg := snapshot.Engine
	// for test
	if config.EngineTerminalAgent == engineCfg.Type && engineCfg.TerminalAgent.Selected == config.AgentClaudeCode {
		slog.Info("Try Translate using claude code")
		cc := engine.NewClaudeCode()
		resCh, err := cc.TranslateStream(context.Background(), engine.Request{
			Prompt:       snapshot.Prompt.Template,
			SystemPrompt: snapshot.Prompt.SystemPrompt,
			Text:         text,
			SourceLang:   sourceLang,
			TargetLang:   targetLang,
		})
		if err != nil {
			return err
		}
		for res := range resCh {
			if res.Text != "" {
				ts.app.Event.Emit("translate", res.Text)
			}
		}
	}
	return nil
}

// ServiceStartup is called when the service starts
func (ts *TranslateService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Store the application instance for later use
	ts.app = application.Get()
	return nil
}

func (ts *TranslateService) ServiceShutdown() error {
	return nil
}
