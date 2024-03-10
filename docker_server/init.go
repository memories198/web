package docker_server

import (
	"context"
	"errors"
	docker "github.com/fsouza/go-dockerclient"
	"time"
)

type Client docker.Client

func PingServer(server string) (*Client, error) {
	cli, err := docker.NewClient(server)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c := make(chan string)
	go func() {
		err = cli.Ping()
		if err == nil {
			c <- ""
		}
	}()
	select {
	case <-ctx.Done():
		return nil, errors.New("连接docker服务器超时")
	case <-c:
	}
	return (*Client)(cli), err
}
