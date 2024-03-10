package docker_server

import (
	"errors"
	docker "github.com/fsouza/go-dockerclient"
	"io"
	"os"
	"strings"
	"time"
)

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

func PullImage(image string, cli *Client, configuration AuthConfiguration) error {
	client := docker.Client(*cli)
	if n := strings.Split(image, ":"); len(n) == 1 {
		image = image + ":latest"
	}

	err := client.PullImage(docker.PullImageOptions{
		Repository: image,
	}, docker.AuthConfiguration(configuration))
	if err != nil {
		return errors.New("拉取镜像失败，没有登录私有镜像或镜像名称错误")

	}
	return nil
}

func PushImage(image string, cli *Client, configuration AuthConfiguration) error {
	client := docker.Client(*cli)
	_, err := client.InspectImage(image)
	if err != nil {
		return err
	}

	err = client.PushImage(docker.PushImageOptions{
		Name: image,
	}, docker.AuthConfiguration(configuration))
	if err != nil {
		return errors.New("未登录到镜像仓库或镜像名称未修改成标准格式")
	}

	return nil
}
