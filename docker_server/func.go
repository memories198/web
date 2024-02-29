package docker_server

import (
	"encoding/json"
	"errors"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"io"
	"os"
	"strings"
	"time"
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

func GetImages(all bool, cli *Client) ([]docker.APIImages, error) {
	client := docker.Client(*cli)
	images, err := client.ListImages(docker.ListImagesOptions{All: all})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func BuildImage(file string, imageName string, cli *Client) (string, error) {
	client := docker.Client(*cli)
	fileStream, err := os.OpenFile(file, os.O_RDONLY, 0440)
	if err != nil {
		return "", err
	}
	defer fileStream.Close()

	err = client.BuildImage(docker.BuildImageOptions{
		Name:         imageName,
		InputStream:  fileStream,
		OutputStream: io.Discard,
	})
	if err != nil {
		return "", err
	}
	image, _ := client.InspectImage(imageName)
	return image.ID[7:], nil
}
func RemoveImage(imageNameOrID string, cli *Client) error {
	client := docker.Client(*cli)
	err := client.RemoveImageExtended(imageNameOrID, docker.RemoveImageOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	return nil
}
func RenameImage(imageTag, newTag string, cli *Client) error {
	client := docker.Client(*cli)
	err := client.TagImage(imageTag, docker.TagImageOptions{
		Repo: newTag,
		//Tag: newTag,//踩过的坑，给tag打标签没有用
	})
	if err != nil {
		return err
	}

	err = client.RemoveImage(imageTag)
	if err != nil {
		return err
	}
	return nil
}

type ImageInfo struct {
	Name        string `json:"名称"`
	Description string `json:"描述"`
	IsOfficial  bool   `json:"是否为官方镜像"`
	StarCount   int    `json:"星标数（受欢迎程度）"`
}

func SearchImagesRepository(name string, cli *Client) (imagesInfo []*ImageInfo, err error) {
	client := docker.Client(*cli)
	images, err := client.SearchImages(name)
	if err != nil {
		return nil, err
	}

	for _, image := range images {
		imageInfo := &ImageInfo{
			Name:        image.Name,
			Description: image.Description,
			IsOfficial:  image.IsOfficial,
			StarCount:   image.StarCount,
		}
		imagesInfo = append(imagesInfo, imageInfo)
	}

	return imagesInfo, nil
}
func AddImageTag(imageName, newTag string, cli *Client) (imageID string, err error) {
	client := docker.Client(*cli)
	err = client.TagImage(imageName, docker.TagImageOptions{
		Repo:  newTag,
		Force: false,
	})
	if err != nil {
		return "", err
	}

	image, _ := client.InspectImage(newTag)
	return image.ID[7:], nil
}

func ExportImages(fileWriter io.Writer, images []string, cli *Client) error {
	client := docker.Client(*cli)
	err := client.ExportImages(docker.ExportImagesOptions{
		Names:             images,
		OutputStream:      fileWriter,
		InactivityTimeout: time.Second * 3600,
	})
	if err != nil {
		return err
	}
	return nil
}
func LoadImages(filePath string, cli *Client) error {
	client := docker.Client(*cli)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	err = client.LoadImage(docker.LoadImageOptions{
		InputStream:  file,
		OutputStream: io.Discard,
		Context:      nil,
	})
	if err != nil {
		return err
	}
	return nil
}

func PullImage(image string, cli *Client) error {
	client := docker.Client(*cli)
	if n := strings.Split(image, ":"); len(n) == 1 {
		image = image + ":latest"
	}

	for _, configuration := range authConfiguration {
		err := client.PullImage(docker.PullImageOptions{
			Repository: image,
		}, configuration)
		if err == nil {
			return nil
		}
	}

	return errors.New("拉取镜像失败，没有登录私有镜像或镜像名称错误")
}

func PushImage(image string, cli *Client) error {
	client := docker.Client(*cli)
	_, err := client.InspectImage(image)
	if err != nil {
		return err
	}

	for _, configuration := range authConfiguration {
		err = client.PushImage(docker.PushImageOptions{
			Name: image,
		}, configuration)
		if err == nil {
			return nil
		}
	}

	return errors.New("未登录到镜像仓库或镜像名称未修改成标准格式")
}

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
