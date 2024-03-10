package docker_server

import (
	docker "github.com/fsouza/go-dockerclient"
)

type AuthConfiguration docker.AuthConfiguration

func CheckRepository(repositoryUsername, repositoryPassword, repository string, cli *Client) error {
	client := docker.Client(*cli)

	_, err := client.AuthCheck(&docker.AuthConfiguration{
		Username:      repositoryUsername,
		Password:      repositoryPassword,
		ServerAddress: repository,
	})
	if err != nil {
		return err
	}

	return nil
}
