package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/natthphong/go-lambda-template/internal/app"
	"github.com/natthphong/go-lambda-template/internal/httpserver"
)

func main() {
	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		fmt.Println("init error:", err)
		os.Exit(1)
	}
	a.RegisterRoutes()

	env := strings.ToUpper(a.Cfg.Env)

	switch env {
	case "", "LOCAL":
		fmt.Println("Running Echo (LOCAL) …")
		if err := httpserver.StartEcho(ctx, a); err != nil {
			fmt.Println("echo error:", err)
			os.Exit(1)
		}
		return

	case "LOCAL_LAMBDA":
		fmt.Println("Running in LOCAL mode — simulate Lambda Function URL request")

		body := map[string]any{
			"username": "tar",
			"password": "secret",
		}
		bodyBytes, _ := json.Marshal(body)

		req := events.LambdaFunctionURLRequest{
			QueryStringParameters: map[string]string{
				"locale": "th",
			},
			RequestContext: events.LambdaFunctionURLRequestContext{
				HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
					Method: "POST",
					Path:   "/api/v1/login",
				},
			},
			Body: string(bodyBytes),
		}

		lh := httpserver.LambdaHandler(a)
		resp, err := lh(ctx, req)
		if err != nil {
			fmt.Println("lambda error:", err)
			os.Exit(1)
		}

		fmt.Println("---- Response ----")
		fmt.Println("StatusCode:", resp.StatusCode)
		fmt.Println(resp.Body)
		return

	default:
		fmt.Println("Running AWS Lambda mode …")
		lambda.Start(httpserver.LambdaHandler(a))
	}
}
