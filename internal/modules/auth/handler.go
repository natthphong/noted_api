package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/natthphong/go-lambda-template/internal/transport"
	"github.com/pkg/errors"
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

func HandleLogin(_ context.Context, p transport.ApiPayload) (interface{}, int, error) {
	var bodyMap map[string]interface{}
	if err := transport.CoerceBodyToMap(p.Body, &bodyMap); err != nil {
		return nil, http.StatusBadRequest, err
	}

	var in LoginRequest
	if err := transport.MapToStruct(bodyMap, &in); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if strings.TrimSpace(in.Username) == "" || strings.TrimSpace(in.Password) == "" {
		return nil, http.StatusBadRequest, errors.New("username/password is required")
	}
	token := "fake-token-" + in.Username
	locale := "en"
	if v, ok := p.Param["locale"]; ok && strings.TrimSpace(v) != "" {
		locale = v
	}

	return LoginResult{User: in.Username, Token: token, Locale: locale}, http.StatusOK, nil
}
