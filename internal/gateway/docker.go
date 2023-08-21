//go:generate mockgen -destination mock_docker.go -package gateway . DockerClient

package gateway

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
)

type DockerClient interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}

type DockerGateway struct {
	DockerClient DockerClient
}

func (d DockerGateway) GetContainersIDsWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	containers, err := d.DockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	var matchingContainersIDs []string
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.HasPrefix(name, "/"+prefix) {
				matchingContainersIDs = append(matchingContainersIDs, container.ID)
				break
			}
		}
	}

	return matchingContainersIDs, nil
}

func (d DockerGateway) ExtractContainerInfo(ctx context.Context, containerID string) (*string, map[string]string, error) {
	containerInspect, err := d.DockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, nil, err
	}

	endpoint := containerInspect.NetworkSettings.IPAddress
	envVars := map[string]string{}
	for _, ev := range containerInspect.Config.Env {
		vars := strings.SplitN(ev, "=", 2)
		envVars[vars[0]] = vars[1]
	}

	return &endpoint, envVars, nil

}
