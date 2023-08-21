package gateway

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestGetObject(t *testing.T) {
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

func TestGetObject_Error(t *testing.T) {
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

func TestGetObject_NotFound(t *testing.T) {
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

func TestPutObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockWrapper := NewMockMinioClientWrapper(ctrl)
	gateway := MinioGateway{
		MinioWrapper: mockWrapper,
		MinioBucket:  "test-bucket",
	}

	objectName := "existing_object"
	data := []byte("Hello, MinIO!")
	body := bytes.NewReader(data)
	objectSize := int64(22)
	mockWrapper.EXPECT().PutObject(
		ctx,
		"test-bucket",
		objectName,
		body,
		objectSize,
		minio.PutObjectOptions{ContentType: "application/json"},
	).Return(minio.UploadInfo{}, nil)

	err := gateway.PutObject(ctx, objectName, body, objectSize)
	assert.NoError(t, err)
}

func TestPutObject_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockWrapper := NewMockMinioClientWrapper(ctrl)
	gateway := MinioGateway{
		MinioWrapper: mockWrapper,
		MinioBucket:  "test-bucket",
	}

	objectName := "existing_object"
	data := []byte("Hello, MinIO!")
	body := bytes.NewReader(data)
	objectSize := int64(22)
	expectedError := assert.AnError
	mockWrapper.EXPECT().PutObject(
		ctx,
		"test-bucket",
		objectName,
		body,
		objectSize,
		minio.PutObjectOptions{ContentType: "application/json"},
	).Return(minio.UploadInfo{}, expectedError)

	err := gateway.PutObject(ctx, objectName, body, objectSize)
	assert.Equal(t, expectedError, err)
}
