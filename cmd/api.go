package main

import (
	"object-storage-gateway/internal/api"

	"github.com/labstack/echo/v4"
)

const days = 90

func main() {
	e := echo.New()

	api.RegisterHandlers(e, api.New())
	e.Logger.Fatal(e.Start(":3000"))
}
