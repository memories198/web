package web_server

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
	"web/config"
	"web/dao"
	docker "web/docker_server"
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
	containers, err := docker.GetAllContainer(all, userClients[username])
	if err != nil {
		return nil, err
	}

	var containersInfo []*Container

	for _, container := range containers {
		containerInfo := &Container{
			Name:        container.Names[0][1:], // 容器至少有一个名字 name:/nginx,[1:0]表示去掉/
			CpuUsage:    docker.GetCpuUsage(&container, userClients[username]),
			MemoryUsage: docker.GetMemoryUsage(&container, userClients[username]),
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
	images, err := docker.GetImages(all, userClients[username])
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
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置 payload
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	id := uuid.New()
	claims["uuid"] = id

	// 设置签名密钥
	tokenString, err := token.SignedString(serverKey)
	if err != nil {
		return err
	}

	c.SetCookie("token", tokenString, config.CookieExpireTime/1000, "/", "", false, true)

	err = dao.MemorySetKey(username+"TokenUUID", id.String(), config.CookieExpireTime/1000)
	if err != nil {
		return err
	}
	return nil
}

func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return serverKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
