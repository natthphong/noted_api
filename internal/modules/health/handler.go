package health

import (
	"context"
	"net/http"
	"time"

	"github.com/natthphong/go-lambda-template/internal/transport"
)

func HandleHealth(_ context.Context, _ transport.ApiPayload) (interface{}, int, error) {
	return map[string]any{
		"status": "ok",
		"uptime": time.Now(),
	}, http.StatusOK, nil
}
