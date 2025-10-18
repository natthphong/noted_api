package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/natthphong/go-lambda-template/internal/app"
	"github.com/natthphong/go-lambda-template/internal/transport"
)

func LambdaHandler(a *app.App) func(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	return func(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		method := strings.ToUpper(req.RequestContext.HTTP.Method)

		// CORS preflight
		if method == http.MethodOptions {
			return jsonResp(http.StatusNoContent, map[string]any{"ok": true, "time": now()})
		}

		// guard methods
		switch method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			return jsonResp(http.StatusMethodNotAllowed, map[string]any{
				"ok": false, "code": "METHOD_NOT_ALLOWED",
				"message": "allowed: GET, POST, PUT, PATCH, DELETE, OPTIONS", "time": now(),
			})
		}

		// Adapt request payload
		var bodyAny any
		if strings.TrimSpace(req.Body) != "" {
			_ = json.Unmarshal([]byte(req.Body), &bodyAny) // ignore errors → handlers validate themselves
		}

		path := req.RequestContext.HTTP.Path
		payload := transport.ApiPayload{
			Path:  path,
			Param: req.QueryStringParameters,
			Body:  bodyAny,
		}

		h, ok := a.Router.Resolve(method, path)
		if !ok {
			return jsonResp(http.StatusNotFound, map[string]any{
				"ok": false, "code": "NOT_FOUND",
				"message": "route not found", "time": now(),
			})
		}

		data, status, err := h(ctx, payload)
		if err != nil {
			return jsonResp(status, map[string]any{"ok": false, "code": "ERROR", "message": err.Error(), "time": now()})
		}
		return jsonResp(status, map[string]any{"ok": true, "data": data, "time": now()})
	}
}

func jsonResp(status int, body map[string]any) (events.LambdaFunctionURLResponse, error) {
	b, _ := json.Marshal(body)
	return events.LambdaFunctionURLResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
		},
		Body:            string(b),
		IsBase64Encoded: false,
	}, nil
}

func now() string { return time.Now().Format(time.RFC3339) }
