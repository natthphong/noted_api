package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResult struct {
	User   string `json:"user"`
	Token  string `json:"token"`
	Locale string `json:"locale,omitempty"`
}

type HandlerFunc func(context.Context, ApiPayload) (interface{}, int, error)

var routes = map[string]map[string]HandlerFunc{
	"POST": {
		"/api/v1/login": handleLogin,
	},
	"GET": {
		"/api/v1/health": handleHealth,
	},
}

func ResolveHandler(method, path string) (HandlerFunc, bool) {
	method = strings.ToUpper(method)
	path = cleanPath(path)

	m, ok := routes[method]
	if !ok {
		return nil, false
	}
	h, ok := m[path]
	return h, ok
}

// GET /api/v1/health
func handleHealth(_ context.Context, _ ApiPayload) (interface{}, int, error) {
	return map[string]any{
		"status": "ok",
		"uptime": "n/a",
	}, http.StatusOK, nil
}

func first(v, def string) string {
	if strings.TrimSpace(v) != "" {
		return v
	}
	return def
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if i := strings.IndexByte(p, '?'); i >= 0 {
		p = p[:i]
	}
	return p
}

func coerceBodyToMap(body any, out *map[string]interface{}) error {
	switch v := body.(type) {
	case nil:
		*out = map[string]interface{}{}
		return nil
	case map[string]interface{}:
		*out = v
		return nil
	case string:
		if strings.TrimSpace(v) == "" {
			*out = map[string]interface{}{}
			return nil
		}
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return fmt.Errorf("body is string but not valid JSON: %w", err)
		}
		*out = m
		return nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("cannot marshal body: %w", err)
		}
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			return fmt.Errorf("cannot unmarshal body to map: %w", err)
		}
		*out = m
		return nil
	}
}

func mapToStruct(m map[string]interface{}, out any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
