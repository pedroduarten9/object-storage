package gateway_integration_docker

import (
	"context"
	"testing"
	"time"

	"object-storage-gateway/internal/gateway"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestGetContainersWithPrefix(t *testing.T) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	dockerGateway := gateway.DockerGateway{DockerClient: cli}

	resp1, _ := cli.ContainerCreate(ctx, &container.Config{
		Image: "minio/minio",
		Env: []string{
			"MINIO_ROOT_USER=admin",
			"MINIO_ROOT_PASSWORD=admin123",
		},
		Cmd: []string{"server", "/tmp/data"},
	}, nil, nil, nil, "minio-1")
	cli.ContainerStart(ctx, resp1.ID, types.ContainerStartOptions{})

	resp2, _ := cli.ContainerCreate(ctx, &container.Config{
		Image: "minio/minio",
		Env: []string{
			"MINIO_ROOT_USER=admin",
			"MINIO_ROOT_PASSWORD=admin123",
		},
		Cmd: []string{"server", "/tmp/data"},
	}, nil, nil, nil, "no-minio-1")
	cli.ContainerStart(ctx, resp2.ID, types.ContainerStartOptions{})

	// Wait for containers to start
	time.Sleep(2 * time.Second)

	defer func() {
		cli.ContainerRemove(ctx, resp1.ID, types.ContainerRemoveOptions{Force: true})
		cli.ContainerRemove(ctx, resp2.ID, types.ContainerRemoveOptions{Force: true})
	}()

	expectedIDs := []string{resp1.ID}
	containersIDs, err := dockerGateway.GetContainersIDsWithPrefix(ctx, "minio")
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, expectedIDs, containersIDs)
	}
}

func TestGetContainersInfo(t *testing.T) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	dockerGateway := gateway.DockerGateway{DockerClient: cli}

	resp1, _ := cli.ContainerCreate(ctx, &container.Config{
		Image: "minio/minio",
		Env: []string{
			"MINIO_ROOT_USER=admin",
			"MINIO_ROOT_PASSWORD=admin123",
		},
		Cmd: []string{"server", "/tmp/data"},
	}, nil, nil, nil, "minio-1")
	cli.ContainerStart(ctx, resp1.ID, types.ContainerStartOptions{})

	// Wait for containers to start
	time.Sleep(2 * time.Second)

	defer func() {
		cli.ContainerRemove(ctx, resp1.ID, types.ContainerRemoveOptions{Force: true})
	}()

	endpoint, envVars, err := dockerGateway.ExtractContainerInfo(ctx, resp1.ID)
	if assert.NoError(t, err) {
		assert.NotNil(t, *endpoint)
		assert.Contains(t, envVars, "MINIO_ROOT_USER")
		assert.Equal(t, "admin", envVars["MINIO_ROOT_USER"])
		assert.Contains(t, envVars, "MINIO_ROOT_PASSWORD")
		assert.Equal(t, "admin123", envVars["MINIO_ROOT_PASSWORD"])
	}
}
