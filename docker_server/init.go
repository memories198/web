package docker_server

import (
	docker "github.com/fsouza/go-dockerclient"
)

type Client docker.Client

func InitUserClient(servers []string) (map[string]*Client, error) {
	clients := map[string]*Client{}

	for _, server := range servers {
		cli, err := docker.NewClient("http://" + server)
		if err != nil {
			return clients, err
		}
		err = cli.Ping()
		if err != nil {
			return clients, err
		}
		clients[server] = (*Client)(cli)
	}

	return clients, nil
}
func AddServerClient(server string) (*Client, error) {
	cli, err := docker.NewClient("http://" + server)
	if err != nil {
		return nil, err
	}
	return (*Client)(cli), nil
}
