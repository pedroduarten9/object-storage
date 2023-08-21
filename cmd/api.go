package main

import (
	"object-storage-gateway/internal/api"
	"object-storage-gateway/internal/domain"
	"object-storage-gateway/internal/gateway"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.uber.org/zap"
)

const days = 90

func main() {
	e := echo.New()

	logger, _ := zap.NewDevelopment()
	e.Use(api.LoggerMiddleware(logger))
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = api.HttpErrorHandler

	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	dockerGateway := gateway.DockerGateway{DockerClient: cli}
	loadBalancer := domain.MinioLoadBalancerImpl{
		DockerGateway:   dockerGateway,
		InstancesPrefix: "object-storage-amazin-object-storage-node-",
	}

	minioBucket := "minio-bucket"
	api.RegisterHandlers(e, api.API{
		LoadBalancer: loadBalancer,
		MinioBucket:  minioBucket,
	})
	e.Logger.Fatal(e.Start(":3000"))
}
