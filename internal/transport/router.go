package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type HandlerFunc func(context.Context, ApiPayload) (interface{}, int, error)

type Router struct {
	routes map[string]map[string]HandlerFunc // method -> path -> handler
}

func NewRouter() *Router {
	return &Router{routes: map[string]map[string]HandlerFunc{}}
}

func (r *Router) Handle(method, path string, h HandlerFunc) {
	method = strings.ToUpper(method)
	path = cleanPath(path)
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = map[string]HandlerFunc{}
	}
	r.routes[method][path] = h
}

func (r *Router) Resolve(method, path string) (HandlerFunc, bool) {
	method = strings.ToUpper(method)
	path = cleanPath(path)
	m, ok := r.routes[method]
	if !ok {
		return nil, false
	}
	h, ok := m[path]
	return h, ok
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

// Helpers to coerce body to map or struct (use in handlers)
func CoerceBodyToMap(body any, out *map[string]interface{}) error {
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

func MapToStruct(m map[string]interface{}, out any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
