package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/natthphong/go-lambda-template/internal/app"
	mdx "github.com/natthphong/go-lambda-template/internal/middleware"
	"github.com/natthphong/go-lambda-template/internal/transport"
)

func StartEcho(ctx context.Context, a *app.App) error {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
	}))
	e.Use(mdx.Logger(a.Logger))

	// One universal handler that adapts Echo to our Router
	e.Any("/*", func(c echo.Context) error {
		method := c.Request().Method
		path := c.Request().URL.Path

		payload := transport.ApiPayload{
			Path:  path,
			Param: map[string]string{},
			Body:  nil,
		}

		// query params into Param (simple, adjust as needed)
		for k, v := range c.QueryParams() {
			if len(v) > 0 {
				payload.Param[k] = v[0]
			}
		}

		// body
		if c.Request().Body != nil {
			var anyBody any
			if err := json.NewDecoder(c.Request().Body).Decode(&anyBody); err == nil {
				payload.Body = anyBody
			}
		}

		h, ok := a.Router.Resolve(method, path)
		if !ok {
			return c.JSON(http.StatusNotFound, map[string]any{
				"ok": false, "code": "NOT_FOUND",
				"message": "route not found", "time": now(),
			})
		}
		data, status, err := h(c.Request().Context(), payload)
		if err != nil {
			return c.JSON(status, map[string]any{
				"ok": false, "code": "ERROR", "message": err.Error(), "time": now(),
			})
		}
		return c.JSON(status, map[string]any{"ok": true, "data": data, "time": now()})
	})

	addr := a.Cfg.Server.Port
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	if addr == "" {
		addr = ":8080"
	}
	fmt.Printf("Listening on %s\n", addr)
	return e.Start(addr)
}
