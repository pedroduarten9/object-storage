package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"object-storage-gateway/internal/domain"
	"object-storage-gateway/internal/gateway"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/object/:uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(id.String())

	mockMinio := gateway.NewMockMinio(ctrl)
	mockLoadBalancer := domain.NewMockMinioLoadBalancer(ctrl)
	minioBucket := "test-bucket"
	s := ServerInterfaceWrapper{
		Handler: API{
			LoadBalancer: mockLoadBalancer,
			MinioBucket:  minioBucket,
		},
	}

	mockLoadBalancer.EXPECT().GetMinioClient(
		c.Request().Context(),
		minioBucket,
		id.String(),
	).Return(mockMinio, nil)
	minioObject := "something"
	mockMinio.EXPECT().GetObject(
		c.Request().Context(),
		id.String(),
	).Return([]byte(minioObject), nil)

	err := s.GetObject(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, minioObject, rec.Body.String())
	}
}

func TestPutObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	e := echo.New()
	minioObject := []byte("something")
	req := httptest.NewRequest(http.MethodPut, "/object/:uuid", bytes.NewReader(minioObject))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(id.String())

	mockMinio := gateway.NewMockMinio(ctrl)
	mockLoadBalancer := domain.NewMockMinioLoadBalancer(ctrl)
	minioBucket := "test-bucket"
	s := ServerInterfaceWrapper{
		Handler: API{
			LoadBalancer: mockLoadBalancer,
			MinioBucket:  minioBucket,
		},
	}

	mockLoadBalancer.EXPECT().GetMinioClient(
		c.Request().Context(),
		minioBucket,
		id.String(),
	).Return(mockMinio, nil)
	mockMinio.EXPECT().PutObject(
		c.Request().Context(),
		id.String(),
		bytes.NewReader(minioObject),
		c.Request().ContentLength,
	).Return(nil)

	err := s.PutObject(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(minioObject), rec.Body.String())
	}
}
