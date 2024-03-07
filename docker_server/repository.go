package docker_server

import (
	"errors"
	docker "github.com/fsouza/go-dockerclient"
)

var authConfiguration = []docker.AuthConfiguration{{}}

func ListRepositories() (repositories []string) {
	for _, configuration := range authConfiguration {
		repositories = append(repositories, configuration.ServerAddress)
	}
	return repositories
}

func LoginRepository(username, password, serverAddress string, cli *Client) error {
	client := docker.Client(*cli)
	for _, configuration := range authConfiguration {
		if username == configuration.Username && password == configuration.Password && serverAddress == configuration.ServerAddress {
			return errors.New("已登录该镜像仓库,无法重复登录")
		}
	}

	_, err := client.AuthCheck(&docker.AuthConfiguration{
		Username:      username,
		Password:      password,
		ServerAddress: serverAddress,
	})
	if err != nil {
		return err
	}

	authConfiguration = append(authConfiguration, docker.AuthConfiguration{
		Username:      username,
		Password:      password,
		ServerAddress: serverAddress,
	})
	return nil
}

func LogoutRepository(repository string) error {
	for i, configuration := range authConfiguration {
		if configuration.ServerAddress == repository {
			if len(authConfiguration) == i-1 {
				authConfiguration = authConfiguration[:i]
			} else {
				authConfiguration = append(authConfiguration[:i], authConfiguration[i+1:]...)
			}
			return nil
		}
	}
	return errors.New("未查询到登录信息，无法注销该镜像仓库")
}
