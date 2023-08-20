package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Act
	api := New()

	//Assert
	assert.Implements(t, (*ServerInterface)(nil), api)
}

func TestGetObject(t *testing.T) {
	// Arrange
	id := uuid.New()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/object/:uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(id.String())
	s := ServerInterfaceWrapper{
		Handler: New(),
	}
	// Act
	err := s.GetObject(c)

	// Assert
	assert.NoError(t, err)
}
func TestPutObject(t *testing.T) {
	// Arrange
	id := uuid.New()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/object/:uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(id.String())
	s := ServerInterfaceWrapper{
		Handler: New(),
	}
	// Act
	err := s.PutObject(c)

	// Assert
	assert.NoError(t, err)
}
