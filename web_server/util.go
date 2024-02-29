package web_server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	docker "web/docker_server"
	"web/user"
)

func getJsonParam(c *gin.Context) func(param string) (value string, err error) {
	jsonData := map[string]interface{}{}
	err := c.BindJSON(&jsonData)

	return func(param string) (string, error) {
		if err != nil {
			return "", err
		}
		value, err := func() (value string, err error) {

			i, exists := jsonData[param]
			if !exists {
				return "", errors.New("缺少" + param + "字段")
			}

			switch i.(type) {
			case string:
				return i.(string), nil
			}

			return "", errors.New(param + "的值不为字符串")
		}()
		if err != nil {
			return "", err
		}
		return value, err
	}
}
func getJsonAnyParam(c *gin.Context) func(param string) (interface{}, error) {
	jsonData := map[string]interface{}{}
	err := c.BindJSON(&jsonData)

	return func(param string) (interface{}, error) {
		if err != nil {
			return nil, err
		}
		value, err := func() (interface{}, error) {

			i, exists := jsonData[param]
			if !exists {
				return nil, errors.New("缺少" + param + "字段")
			}

			return i, nil
		}()
		if err != nil {
			return nil, err
		}
		return value, err
	}
}

func unixTimeToReadableTime(t int64) string {
	timer := time.Unix(t, 0)
	return timer.Format("2006-01-02 15:04:05 MST")
}

func bytesToMB(size int64) string {
	return strconv.FormatInt(size/1000/1000, 10) + "MB"
}

func containers(all bool, username, server string) ([]*Container, error) {
	containers, err := docker.GetAllContainer(all, userClients[username][server])
	if err != nil {
		return nil, err
	}

	var containersInfo []*Container

	for _, container := range containers {
		containerInfo := &Container{
			Name:        container.Names[0][1:], // 容器至少有一个名字 name:/nginx,[1:0]表示去掉/
			CpuUsage:    docker.GetCpuUsage(&container, userClients[username][server]),
			MemoryUsage: docker.GetMemoryUsage(&container, userClients[username][server]),
			Image:       container.Image,
			ID:          container.ID,
			Status:      container.Status,
			State:       container.State,
			Ports:       container.Ports,
			Mounts:      container.Mounts,
		}
		// 将容器信息添加到切片中
		containersInfo = append(containersInfo, containerInfo)
	}
	return containersInfo, nil
}
func images(all bool, username, server string) ([]Image, error) {
	images, err := docker.GetImages(all, userClients[username][server])
	if err != nil {
		return nil, err
	}

	var imagesInfo []Image

	for _, image := range images {
		imageInfo := Image{
			ID:      image.ID[7:],
			Tags:    image.RepoTags,
			Created: unixTimeToReadableTime(image.Created),
			Size:    bytesToMB(image.Size),
		}
		if image.ParentID != "" {
			imageInfo.ParentID = image.ParentID[7:]
		}

		imagesInfo = append(imagesInfo, imageInfo)
	}
	return imagesInfo, nil
}

func setCookie(c *gin.Context, username string) error {
	cookie, err := user.MemoryLPop(username)
	if err == nil {
		err = user.MemoryDelKey(cookie)
		if err != nil {
			return err
		}
	}

	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	c.SetCookie("username", now, cookieOutTime, "/", "", false, true)
	err = user.MemorySetKey(now, username, cookieOutTime)
	if err != nil {
		return err
	}
	err = user.MemoryLPush(username, now)
	if err != nil {
		return err
	}
	return nil
}
