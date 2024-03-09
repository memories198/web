package docker_server

import (
	"context"
	"errors"
	docker "github.com/fsouza/go-dockerclient"
	"sync"
	"time"
)

type Client docker.Client

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

//	AddServerClient func InitUserClient(servers []string) (clients map[string]*Client, errs []string) {
//		clients = map[string]*Client{}
//
//		for _, server := range servers {
//			wg.Add(1)
//			go func(server string) {
//				defer wg.Done()
//				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//				defer cancel()
//				c := make(chan string)
//				go func() {
//					cli, err := docker.NewClient("http://" + server)
//					if err != nil {
//						c <- server + "的客户端初始化失败:" + err.Error()
//						return
//					}
//					err = cli.Ping()
//					if err != nil {
//						c <- server + "无法连接:" + err.Error()
//						return
//					}
//					clients[server] = (*Client)(cli)
//
//					c <- ""
//				}()
//				select {
//				case <-ctx.Done():
//					errs = append(errs, server+"的客户端初始化失败或连接超时")
//					return
//
//				case message := <-c:
//					if message != "" {
//						errs = append(errs, message)
//					}
//					return
//				}
//			}(server)
//		}
//		wg.Wait()
//		return
//	}
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
