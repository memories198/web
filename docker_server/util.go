package docker_server

import (
	docker "github.com/fsouza/go-dockerclient"
	"log"
)

func getContainerStats(container *docker.APIContainers, client *docker.Client) *docker.Stats {
	var containerStatsChannel = make(chan *docker.Stats)
	done := make(chan bool)
	go func() {
		// 调用Stats方法获取容器的统计数据，设置Stream为false以禁用流模式
		err := client.Stats(docker.StatsOptions{
			ID:     container.ID,
			Stats:  containerStatsChannel,
			Stream: true, // 设置为false以禁用流模式
			Done:   nil,
		})
		if err != nil {
			log.Println("获取容器统计信息失败:", err)
		}
		// 通知主Goroutine可以继续执行并清理资源
		done <- true
	}()

	return <-containerStatsChannel
}
func SearchImagesLocal(name string, cli *Client) bool {
	images, err := GetImages(false, cli)
	if err != nil {
		return false
	}
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == name {
				return true
			}
		}
	}
	return false
}
