package gateway_integration_minio

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"object-storage-gateway/internal/gateway"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

const (
	bucketName = "test-bucket"
)

var (
	testMinioClient *minio.Client
	testServer      *httptest.Server
	containerID     string
)

func setupMinioServer() {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	user := "admin"
	pass := "admin123"
	cli.ImagePull(ctx, "minio/minio", types.ImagePullOptions{})

	portBinding := map[nat.Port][]nat.PortBinding{
		"9000/tcp": {{HostIP: "0,0,0,0", HostPort: "9000"}},
	}
	resp, _ := cli.ContainerCreate(ctx, &container.Config{
		Image: "minio/minio",
		Env: []string{
			"MINIO_ROOT_USER=" + user,
			"MINIO_ROOT_PASSWORD=" + pass,
		},
		ExposedPorts: nat.PortSet{"9000/tcp": struct{}{}},
		Cmd:          []string{"server", "/tmp/data"},
	}, &container.HostConfig{
		PortBindings: portBinding,
	}, nil, nil, "test-minio")
	containerID = resp.ID

	cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})

	// Wait for containers to start
	time.Sleep(2 * time.Second)

	testMinioClient, _ = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(user, pass, ""),
		Secure: false,
	})
}

func teardown() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return
	}

	err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return
	}

	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
	if err != nil {
		return
	}
}

func TestMain(m *testing.M) {
	// Setup
	setupMinioServer()

	// Run tests
	exitCode := m.Run()

	// Teardown
	teardown()

	os.Exit(exitCode)
}

func TestPutGetObject(t *testing.T) {
	ctx := context.Background()
	minioGateway := gateway.MinioGateway{
		MinioWrapper: gateway.MinioWrapper{
			MinioClient: testMinioClient,
		},
		MinioBucket: bucketName,
	}

	objectName := "example.json"

	data := []byte("Hello, MinIO!")
	err := minioGateway.PutObject(ctx, objectName, bytes.NewReader(data), int64(len(data)))
	assert.NoError(t, err)

	retrievedData, err := minioGateway.GetObject(ctx, objectName)
	if assert.NoError(t, err) {
		assert.Equal(t, data, retrievedData)
	}
}

func TestGetObject_NotFound(t *testing.T) {
	ctx := context.Background()
	minioGateway := gateway.MinioGateway{
		MinioWrapper: gateway.MinioWrapper{
			MinioClient: testMinioClient,
		},
		MinioBucket: bucketName,
	}

	objectName := "not_found.json"
	retrievedData, err := minioGateway.GetObject(ctx, objectName)
	assert.Nil(t, retrievedData)
	assert.Equal(t, gateway.NotFoundError{Msg: fmt.Sprintf("object %s not found", objectName)}, err)
}
