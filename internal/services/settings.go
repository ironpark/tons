package services

import (
	"context"

	"github.com/ironpark/tons/internal/config"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type SettingService struct {
	cfg *config.Config
	app *application.App
}

func NewSettingService() (*SettingService, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return &SettingService{
		cfg: cfg,
	}, nil
}

func (ss *SettingService) GetCurrentConfig() *config.Config {
	return ss.cfg.Snapshot()
}

func (ss *SettingService) UpdateGeneralConfig(general config.GeneralConfig) error {
	ss.cfg.SetGeneral(general)
	return ss.cfg.Save()
}

func (ss *SettingService) UpdateEngineConfig(engine config.EngineConfig) error {
	ss.cfg.SetEngine(engine)
	return ss.cfg.Save()
}

func (ss *SettingService) UpdatePromptConfig(prompt config.PromptConfig) error {
	ss.cfg.SetPromptConfig(prompt)
	return ss.cfg.Save()
}

// ServiceStartup is called when the service starts
func (ss *SettingService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Store the application instance for later use
	ss.app = application.Get()
	return nil
}

func (u *SettingService) ServiceShutdown() error {
	return nil
}
