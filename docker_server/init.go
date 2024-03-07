package docker_server

import (
	"context"
	"errors"
	docker "github.com/fsouza/go-dockerclient"
	"sync"
	"time"
)

type Client docker.Client

func InitUserClient(servers []string) (clients map[string]*Client, errs []string) {
	clients = map[string]*Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, server := range servers {
		wg.Add(1)
		go func(ctx context.Context, se string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				mu.Lock()
				defer mu.Unlock()
				errs = append(errs, se+"的客户端初始化或连接超时")
				return
			case <-func() <-chan string {
				c := make(chan string)

				cli, err := docker.NewClient("http://" + se)
				if err != nil {
					mu.Lock()
					defer mu.Unlock()
					errs = append(errs, se+"的客户端初始化失败:"+err.Error())
					close(c)
					return c
				}
				err = cli.Ping()
				if err != nil {
					mu.Lock()
					defer mu.Unlock()
					errs = append(errs, se+"无法连接:"+err.Error())
					close(c)
					return c
				}
				mu.Lock()
				defer mu.Unlock()
				clients[se] = (*Client)(cli)
				close(c)
				return c
			}():
			}

		}(ctx, server)
	}
	wg.Wait()
	return
}
func AddServerClient(server string) (*Client, error) {
	var (
		err error
		cli *docker.Client
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		err = errors.New(server + "的客户端初始化或连接超时")
	case <-func() chan string {
		c := make(chan string)

		cli, err = docker.NewClient("http://" + server)
		if err != nil {
			err = errors.New(server + "无法添加:" + err.Error())
		}
		close(c)
		return c
	}():
	}

	return (*Client)(cli), err
}
