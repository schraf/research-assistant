package service

import (
	"context"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func StartServer(ctx context.Context) error {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}

	if err := funcframework.StartHostPort(hostname, port); err != nil {
		return err
	}

	return nil
}
