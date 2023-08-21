package gateway

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestGetContainersIDsWithPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	prefix := "some-prefix"
	mockDockerClient := NewMockDockerClient(ctrl)
	dockerGateway := DockerGateway{
		DockerClient: mockDockerClient,
	}

	containers := []types.Container{
		{
			ID:    "some-id1",
			Names: []string{"/some-prefixekosnfwnfo"},
		},
		{
			ID:    "some-id2",
			Names: []string{"/some-prefix-som"},
		},
		{
			ID:    "some-id3",
			Names: []string{"/some-prefix-nothing"},
		},
		{
			ID:    "some-id4",
			Names: []string{"/no-prefix"},
		},
	}
	mockDockerClient.EXPECT().ContainerList(
		ctx,
		types.ContainerListOptions{},
	).Return(containers, nil)
	expectedContainersIDs := []string{"some-id1", "some-id2", "some-id3"}
	containersIDs, err := dockerGateway.GetContainersIDsWithPrefix(ctx, prefix)

	if assert.NoError(t, err) {
		assert.ElementsMatch(t, expectedContainersIDs, containersIDs)
	}
}

func TestGetContainersIDsWithPrefix_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	prefix := "some-prefix"
	mockDockerClient := NewMockDockerClient(ctrl)
	dockerGateway := DockerGateway{
		DockerClient: mockDockerClient,
	}

	containers := []types.Container{}
	mockDockerClient.EXPECT().ContainerList(
		ctx,
		types.ContainerListOptions{},
	).Return(containers, nil)
	containersIDs, err := dockerGateway.GetContainersIDsWithPrefix(ctx, prefix)

	if assert.NoError(t, err) {
		assert.Empty(t, containersIDs)
	}
}

func TestGetContainersIDsWithPrefix_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	prefix := "some-prefix"
	mockDockerClient := NewMockDockerClient(ctrl)
	dockerGateway := DockerGateway{
		DockerClient: mockDockerClient,
	}

	expectedError := assert.AnError
	mockDockerClient.EXPECT().ContainerList(
		ctx,
		types.ContainerListOptions{},
	).Return(nil, expectedError)
	containersIDs, err := dockerGateway.GetContainersIDsWithPrefix(ctx, prefix)

	if assert.Equal(t, expectedError, err) {
		assert.Nil(t, containersIDs)
	}
}

func TestExtractContainerInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	containerID := "some-id"
	expectedEndpoint := "some-endpoint"
	expectedEnvVars := map[string]string{
		"env-var-1": "value1",
		"env-var-2": "value2",
	}
	mockDockerClient := NewMockDockerClient(ctrl)
	dockerGateway := DockerGateway{
		DockerClient: mockDockerClient,
	}

	container := types.ContainerJSON{
		NetworkSettings: &types.NetworkSettings{
			DefaultNetworkSettings: types.DefaultNetworkSettings{
				IPAddress: expectedEndpoint,
			},
		},
		Config: &container.Config{
			Env: []string{"env-var-1=value1", "env-var-2=value2"},
		},
	}

	mockDockerClient.EXPECT().ContainerInspect(
		ctx,
		containerID,
	).Return(container, nil)

	endpoint, envVars, err := dockerGateway.ExtractContainerInfo(ctx, containerID)
	if assert.NoError(t, err) {
		assert.Equal(t, expectedEndpoint, *endpoint)
		assert.Equal(t, expectedEnvVars, envVars)
	}
}

func TestExtractContainerInfo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	containerID := "some-id"
	mockDockerClient := NewMockDockerClient(ctrl)
	dockerGateway := DockerGateway{
		DockerClient: mockDockerClient,
	}

	expectedError := assert.AnError
	mockDockerClient.EXPECT().ContainerInspect(
		ctx,
		containerID,
	).Return(types.ContainerJSON{}, expectedError)

	endpoint, envVars, err := dockerGateway.ExtractContainerInfo(ctx, containerID)
	if assert.Equal(t, expectedError, err) {
		assert.Nil(t, endpoint)
		assert.Nil(t, envVars)
	}
}
