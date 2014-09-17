package container

import (
	"github.com/fsouza/go-dockerclient"
)

var dockerClient *docker.Client

func GetDockerClient() (*docker.Client, error) {
	if dockerClient != nil {
		return dockerClient, nil
	}
	endpoint := "unix:///var/run/docker.sock"
	dockerClient, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	return dockerClient, nil
}
