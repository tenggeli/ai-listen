package app

import (
	"fmt"

	"ai-listen/backend/internal/config"
	"ai-listen/backend/internal/router"
	"go.uber.org/zap"
)

type App struct {
	cfg    config.Config
	logger *zap.Logger
}

func New() (*App, error) {
	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	return &App{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (a *App) Run() error {
	defer func() {
		_ = a.logger.Sync()
	}()

	engine := router.New(a.logger)
	a.logger.Info("listen backend starting", zap.String("addr", a.cfg.HTTPAddr))

	return engine.Run(a.cfg.HTTPAddr)
}
