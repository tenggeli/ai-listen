package app

import (
	"fmt"

	"ai-listen/backend/internal/config"
	"ai-listen/backend/internal/router"
	"ai-listen/backend/internal/store"
	"go.uber.org/zap"
)

type App struct {
	cfg    config.Config
	logger *zap.Logger
	db     *store.MySQLStore
}

func New() (*App, error) {
	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	app := &App{
		cfg:    cfg,
		logger: logger,
	}

	if cfg.MySQLDSN != "" {
		db, err := store.NewMySQLStore(cfg.MySQLDSN)
		if err != nil {
			return nil, fmt.Errorf("init mysql store: %w", err)
		}
		store.SetDefault(db)
		app.db = db
	}

	return app, nil
}

func (a *App) Run() error {
	defer func() {
		if a.db != nil {
			_ = a.db.Close()
		}
		_ = a.logger.Sync()
	}()

	engine := router.New(a.logger)
	a.logger.Info("listen backend starting", zap.String("addr", a.cfg.HTTPAddr))

	return engine.Run(a.cfg.HTTPAddr)
}
