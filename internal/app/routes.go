package app

import (
	"github.com/natthphong/go-lambda-template/modules/auth"
	"github.com/natthphong/go-lambda-template/modules/health"
)

func (a *App) RegisterRoutes() {
	// v1
	a.Router.Handle("GET", "/api/v1/health", health.HandleHealth)
	a.Router.Handle("POST", "/api/v1/login", auth.HandleLogin)

	// later (v2)
	// a.Router.Handle("POST", "/api/v2/login", auth.HandleLoginV2)
}
