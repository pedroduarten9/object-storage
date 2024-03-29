package api

import (
	"fmt"
	"net/http"
	"object-storage-gateway/internal/gateway"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.Time("Start time", v.StartTime),
				zap.Duration("Latency", v.Latency),
				zap.String("URI", v.URI),
				zap.Int("Status", v.Status),
			)

			return nil
		},
	})
}

func HttpErrorHandler(err error, ctx echo.Context) {
	switch e := err.(type) {
	case gateway.NotFoundError:
		ctx.JSON(http.StatusNotFound, Error{Message: err.Error()})
	default:
		fmt.Print(e)
		ctx.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
}
