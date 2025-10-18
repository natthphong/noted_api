package app

import (
	"context"

	"github.com/natthphong/go-lambda-template/config"
	"github.com/natthphong/go-lambda-template/internal/logz"
	"github.com/natthphong/go-lambda-template/internal/transport"
	"go.uber.org/zap"
)

type App struct {
	Cfg    *config.Config
	Logger *zap.Logger
	Router *transport.Router
}

func New(ctx context.Context) (*App, error) {
	config.InitTimeZone()
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}
	logz.Init(cfg.Log.Level, cfg.Server.Name)
	logger := zap.L()
	r := transport.NewRouter()
	return &App{Cfg: cfg, Logger: logger, Router: r}, nil
}
