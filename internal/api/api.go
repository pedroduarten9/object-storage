//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config generation-config.yaml ../../object-storage-gateway.yaml
package api

import "github.com/labstack/echo/v4"

type API struct{}

func New() ServerInterface {
	return &API{}
}

func (a API) GetObject(ctx echo.Context, uuid UuidPath) error {
	return nil
}

func (a API) PutObject(ctx echo.Context, uuid UuidPath) error {
	return nil
}
