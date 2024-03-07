package docker_server

import (
	"encoding/json"
	"errors"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"strings"
)

func GetAllContainer(all bool, cli *Client) ([]docker.APIContainers, error) {
	client := docker.Client(*cli)
	opts := docker.ListContainersOptions{All: all} // 设置 All: false 以仅列出运行中的容器，设置为 true 则包括停止的容器
	containers, err := client.ListContainers(opts)
	if err != nil {
		return nil, err
	}
	return containers, nil
}
func CreateContainer(configBytes []byte, cli *Client) (ID string, err error) {
	client := docker.Client(*cli)
	var config docker.CreateContainerOptions
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return "", err
	}

	image := config.Config.Image
	if n := strings.Split(image, ":"); len(n) == 1 {
		image = image + ":latest"
	}
	exist := SearchImagesLocal(image, cli)
	pullOK := false
	if exist == false {
		for _, configuration := range authConfiguration {
			err := client.PullImage(docker.PullImageOptions{
				Repository: image,
			}, configuration)
			if err == nil {
				pullOK = true
				break
			}
		}
	}
	if exist == false && !pullOK {
		return "", errors.New("拉取镜像失败")
	}

	container, err := client.CreateContainer(config)
	if err != nil {
		return "", err
	}
	return container.ID, nil
}

func StartContainer(ID string, cli *Client) error {
	client := docker.Client(*cli)
	err := client.StartContainer(ID, nil)
	if err != nil {
		return err
	}
	return nil
}

func StopContainer(ID string, cli *Client) error {
	client := docker.Client(*cli)
	err := client.StopContainer(ID, 7)
	if err != nil {
		return err
	}
	return nil
}

func RemoveContainer(ID string, cli *Client) error {
	client := docker.Client(*cli)
	err := client.RemoveContainer(docker.RemoveContainerOptions{ID: ID})
	if err != nil {
		return err
	}
	return nil
}

func KillContainer(ID string, cli *Client) error {
	client := docker.Client(*cli)
	err := client.KillContainer(docker.KillContainerOptions{ID: ID})
	if err != nil {
		return err
	}
	return nil
}

func GetCpuUsage(container *docker.APIContainers, cli *Client) (cpu string) {
	client := docker.Client(*cli)
	// 计算CPU使用百分比
	stats := getContainerStats(container, &client)

	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	//容器从现在到上一个时间点使用的cpu总时间
	systemDelta := float64(stats.CPUStats.SystemCPUUsage - stats.PreCPUStats.SystemCPUUsage)
	//系统上所有进程从现在到上一个时间点使用的cpu总时间
	cpuPercent := cpuDelta / systemDelta * 100.0

	cpu = fmt.Sprintf("%.1f%%", cpuPercent)

	return
}

func GetMemoryUsage(container *docker.APIContainers, cli *Client) string {
	client := docker.Client(*cli)
	stats := getContainerStats(container, &client)
	// 提取内存使用情况
	memoryUsage := stats.MemoryStats.Usage

	return fmt.Sprintf("%vMiB", memoryUsage/1e6)
}
