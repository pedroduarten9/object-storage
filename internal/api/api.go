//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config generation-config.yaml ../../object-storage-gateway.yaml
package api

import (
	"bytes"
	"io"
	"net/http"
	"object-storage-gateway/internal/domain"

	"github.com/labstack/echo/v4"
)

type API struct {
	MinioBucket  string
	LoadBalancer domain.MinioLoadBalancer
}

func (a API) GetObject(ctx echo.Context, uuid UuidPath) error {
	minioClient, err := a.LoadBalancer.GetMinioClient(ctx.Request().Context(), a.MinioBucket, string(uuid))
	if err != nil {
		return err
	}
	object, err := minioClient.GetObject(ctx.Request().Context(), uuid)
	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, string(object))
}

func (a API) PutObject(ctx echo.Context, uuid UuidPath) error { // Read the Body content
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, ctx.Request().Body)
	if err != nil {
		return err
	}

	minioClient, err := a.LoadBalancer.GetMinioClient(ctx.Request().Context(), a.MinioBucket, string(uuid))
	if err != nil {
		return err
	}

	err = minioClient.PutObject(ctx.Request().Context(), string(uuid), bytes.NewReader(buf.Bytes()), ctx.Request().ContentLength)
	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, string(buf.Bytes()))
}
