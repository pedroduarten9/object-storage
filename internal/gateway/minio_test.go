package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestMinioGateway_GetExistingObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockWrapper := NewMockMinioClientWrapper(ctrl)
	mockObject := NewMockMinioObject(ctrl)
	gateway := MinioGateway{
		MinioWrapper: mockWrapper,
		MinioBucket:  "test-bucket",
	}

	objectName := "existing_object"
	expectedObject := []byte{}
	mockWrapper.EXPECT().GetObject(
		ctx,
		"test-bucket",
		objectName,
		minio.GetObjectOptions{},
	).Return(mockObject, nil).Times(1)
	mockObject.EXPECT().Read(
		gomock.Any(),
	).DoAndReturn(func(p []byte) (int, error) {
		n := copy(p, expectedObject)
		if n == 0 {
			return 0, io.EOF
		}
		return n, nil
	})

	data, err := gateway.GetObject(ctx, objectName)

	if assert.NoError(t, err) {
		assert.Equal(t, expectedObject, data)
	}
}

func TestMinioGateway_ErrorGetObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockWrapper := NewMockMinioClientWrapper(ctrl)
	gateway := MinioGateway{
		MinioWrapper: mockWrapper,
		MinioBucket:  "test-bucket",
	}

	expectedError := assert.AnError
	objectName := "non_existing_object"
	mockWrapper.EXPECT().GetObject(
		ctx,
		"test-bucket",
		objectName,
		minio.GetObjectOptions{},
	).Return(nil, expectedError)
	data, err := gateway.GetObject(ctx, objectName)
	if assert.Equal(t, expectedError, err) {
		assert.Nil(t, data)
	}
}

func TestMinioGateway_GetNonExistingObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockWrapper := NewMockMinioClientWrapper(ctrl)
	mockObject := NewMockMinioObject(ctrl)
	gateway := MinioGateway{
		MinioWrapper: mockWrapper,
		MinioBucket:  "test-bucket",
	}

	objectName := "non_existing_object"
	mockWrapper.EXPECT().GetObject(
		ctx,
		"test-bucket",
		objectName,
		minio.GetObjectOptions{},
	).Return(mockObject, nil)
	mockObject.EXPECT().Read(
		gomock.Any(),
	).DoAndReturn(func(p []byte) (int, error) {
		return 0, errors.New("object not found")
	})
	data, err := gateway.GetObject(ctx, objectName)
	if assert.Equal(t, &NotFoundError{fmt.Sprintf("object %s not found", objectName)}, err) {
		assert.Nil(t, data)
	}
}
