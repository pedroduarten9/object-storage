package main

import (
	"object-storage-gateway/internal/api"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

const days = 90

func main() {
	e := echo.New()

	logger, _ := zap.NewDevelopment()
	e.Use(api.LoggerMiddleware(logger))
	e.Use(api.RequestIDMiddleware())
	e.Use(api.AuthenticationMiddleware())
	e.Use(middleware.Recover())

	api.RegisterHandlers(e, api.New())
	e.Logger.Fatal(e.Start(":3000"))
}
