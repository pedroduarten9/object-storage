//go:generate mockgen -destination mock_minio.go -package gateway . MinioClient,MinioObject,MinioClientWrapper

package gateway

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioObject interface {
	Read(b []byte) (n int, err error)
}

var _ MinioObject = (*minio.Object)(nil)

type MinioClient interface {
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
}

var _ MinioClient = (*minio.Client)(nil)

type MinioClientWrapper interface {
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (MinioObject, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
}

var _ MinioClientWrapper = (*MinioWrapper)(nil)

type MinioWrapper struct {
	MinioClient *minio.Client
}

func (mw MinioWrapper) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (MinioObject, error) {
	return mw.MinioClient.GetObject(ctx, bucketName, objectName, opts)
}

func (mw MinioWrapper) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	return mw.MinioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

type MinioGateway struct {
	MinioWrapper MinioClientWrapper
	MinioBucket  string
}

func (m MinioGateway) GetObject(ctx context.Context, objectName string) ([]byte, error) {
	object, err := m.MinioWrapper.GetObject(
		ctx,
		m.MinioBucket,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, object)
	if err != nil {
		return nil, &NotFoundError{fmt.Sprintf("object %s not found", objectName)}
	}

	return buf.Bytes(), nil
}

func (m MinioGateway) PutObject(ctx context.Context, objectName string, body *bytes.Reader, objectSize int64) error {
	_, err := m.MinioWrapper.PutObject(
		ctx,
		m.MinioBucket,
		objectName,
		body,
		objectSize,
		minio.PutObjectOptions{ContentType: "application/json"},
	)
	if err != nil {
		return err
	}

	return nil
}
