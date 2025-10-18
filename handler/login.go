package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// POST /api/v1/login
func handleLogin(_ context.Context, p ApiPayload) (interface{}, int, error) {
	var bodyMap map[string]interface{}
	if err := coerceBodyToMap(p.Body, &bodyMap); err != nil {
		return nil, http.StatusBadRequest, err
	}

	var in LoginRequest
	if err := mapToStruct(bodyMap, &in); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if strings.TrimSpace(in.Username) == "" || strings.TrimSpace(in.Password) == "" {
		return nil, http.StatusBadRequest, errors.New("username/password is required")
	}

	token := "fake-token-" + in.Username
	return LoginResult{
		User:   in.Username,
		Token:  token,
		Locale: first(p.Param["locale"], "en"),
	}, http.StatusOK, nil
}
