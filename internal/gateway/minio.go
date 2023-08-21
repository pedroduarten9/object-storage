//go:generate mockgen -destination mock_minio.go -package gateway . MinioObject,MinioClientWrapper,Minio

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

type MinioClientWrapper interface {
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (MinioObject, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
	CreateBucketIfNotExists(ctx context.Context, bucketName string) error
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

func (mw MinioWrapper) CreateBucketIfNotExists(ctx context.Context, bucketName string) error {
	exists, err := mw.MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		mw.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}
	return nil
}

type Minio interface {
	GetObject(ctx context.Context, objectName string) ([]byte, error)
	PutObject(ctx context.Context, objectName string, body *bytes.Reader, objectSize int64) error
}

type MinioGateway struct {
	MinioWrapper MinioClientWrapper
	MinioBucket  string
}

func (m MinioGateway) GetObject(ctx context.Context, objectName string) ([]byte, error) {
	err := m.MinioWrapper.CreateBucketIfNotExists(ctx, m.MinioBucket)
	if err != nil {
		return nil, err
	}

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
		return nil, NotFoundError{fmt.Sprintf("object %s not found", objectName)}
	}

	return buf.Bytes(), nil
}

func (m MinioGateway) PutObject(ctx context.Context, objectName string, body *bytes.Reader, objectSize int64) error {
	err := m.MinioWrapper.CreateBucketIfNotExists(ctx, m.MinioBucket)
	if err != nil {
		return err
	}

	_, err = m.MinioWrapper.PutObject(
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

func (m MinioGateway) CreateBucketIfNotExists(ctx context.Context, bucketName string) error {
	return m.MinioWrapper.CreateBucketIfNotExists(ctx, bucketName)
}
