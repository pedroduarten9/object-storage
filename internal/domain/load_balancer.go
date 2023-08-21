//go:generate mockgen -destination mock_load_balancer.go -package domain . MinioLoadBalancer
package domain

import (
	"context"
	"object-storage-gateway/internal/gateway"

	"github.com/dgryski/go-farm"
	jump "github.com/dgryski/go-jump"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	minioAccessKey = "MINIO_ACCESS_KEY"
	minioSecretKey = "MINIO_SECRET_KEY"
)

type MinioLoadBalancer interface {
	GetMinioClient(ctx context.Context, minioBucket string, uuid string) (gateway.Minio, error)
}

type MinioLoadBalancerImpl struct {
	DockerGateway   gateway.DockerGateway
	InstancesPrefix string
}

func (l MinioLoadBalancerImpl) GetMinioClient(ctx context.Context, minioBucket string, uuid string) (gateway.Minio, error) {
	containersIDs, err := l.DockerGateway.GetContainersIDsWithPrefix(ctx, l.InstancesPrefix)
	if err != nil {
		return nil, err
	}

	node := getContainer(containersIDs, uuid)
	endpoint, envVars, err := l.DockerGateway.ExtractContainerInfo(ctx, node)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(*endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(envVars[minioAccessKey], envVars[minioSecretKey], ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return gateway.MinioGateway{
		MinioWrapper: gateway.MinioWrapper{
			MinioClient: minioClient,
		},
		MinioBucket: minioBucket,
	}, nil
}

func Sum64(data []byte) uint64 {
	return farm.Hash64(data)
}

func getContainer(containersIDs []string, uuid string) string {
	containerIdx := jump.Hash(farm.Hash64([]byte(uuid)), len(containersIDs))
	return containersIDs[containerIdx]
}
