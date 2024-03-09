package web_server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	docker "web/docker_server"
)

type Container struct {
	Name        string      `json:"Name"`
	CpuUsage    string      `json:"CpuUsage"`
	MemoryUsage string      `json:"MemoryUsage"`
	Image       string      `json:"Image"`
	ID          string      `json:"ID"`
	Status      string      `json:"Status"`
	State       string      `json:"State"`
	Ports       interface{} `json:"Ports"`
	Mounts      interface{} `json:"Mounts"`
}

func containersList(c *gin.Context) {
	username, _ := c.Get("username")
	server := c.Query("server")

	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	allContainers, err := containers(false, username.(string), server)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取所有容器信息失败",
			"error":   err.Error(),
		})
		return
	}

	if allContainers == nil {
		c.JSON(200, gin.H{
			"message": "没有容器正在运行",
		})
		return
	}

	marshal, err := json.Marshal(allContainers)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "查看所有容器失败",
			"error":   err.Error(),
		})
		return
	}

	// 发送JSON数据
	c.Data(http.StatusOK, "application/json", marshal)
}
func containersListAll(c *gin.Context) {
	username, _ := c.Get("username")
	server := c.Query("server")

	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	allContainers, err := containers(true, username.(string), server)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "查看所有容器失败",
			"error":   err.Error(),
		})
	}

	if allContainers == nil {
		c.JSON(200, gin.H{
			"message": "没有容器正在运行",
		})
		return
	}

	marshal, err := json.Marshal(allContainers)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"message": "将查找到的容器信息转换为json失败",
			"error":   err.Error(),
		})
		return
	}

	// 发送JSON数据
	c.Data(http.StatusOK, "application/json", marshal)
}
func containerCreate(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "json解析失败",
			"error":   err.Error(),
		})
		return
	}
	username, _ := c.Get("username")

	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	ID, err := docker.CreateContainer(body, userClients[username.(string)])

	if err != nil {
		c.JSON(400, gin.H{
			"message": "容器创建失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message":     "容器创建成功",
		"containerID": ID,
	})
}

func containerStart(c *gin.Context) {
	data := getJsonParam(c)
	IDOrName, err := data("IDOrName")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "JSON解析出错",
			"error":   err.Error(),
		})
		return
	}
	username, _ := c.Get("username")
	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.StartContainer(IDOrName, userClients[username.(string)])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":           "启动容器失败",
			"containerIDOrName": IDOrName,
			"error":             err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":           "启动容器成功",
		"containerIDOrName": IDOrName,
	})
}
func containerStop(c *gin.Context) {
	data := getJsonParam(c)
	IDOrName, err := data("IDOrName")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "JSON解析出错",
			"error":   err.Error(),
		})
		return
	}
	username, _ := c.Get("username")
	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.StopContainer(IDOrName, userClients[username.(string)])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":           "容器暂停失败",
			"containerIDOrName": IDOrName,
			"error":             err.Error(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message":           "容器暂停成功",
		"containerIDOrName": IDOrName,
	})
}
func containerKill(c *gin.Context) {
	data := getJsonParam(c)
	IDOrName, err := data("IDOrName")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "JSON解析出错",
			"error":   err.Error(),
		})
		return
	}
	username, _ := c.Get("username")
	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.KillContainer(IDOrName, userClients[username.(string)])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":           "容器强制暂停失败",
			"containerIDOrName": IDOrName,
			"error":             err.Error(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message":           "容器强制暂停成功",
		"containerIDOrName": IDOrName,
	})
}
func containerRemove(c *gin.Context) {
	data := getJsonParam(c)
	IDOrName, err := data("IDOrName")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "JSON解析出错",
			"error":   err.Error(),
		})
		return
	}
	username, _ := c.Get("username")
	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.RemoveContainer(IDOrName, userClients[username.(string)])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             err.Error(),
			"containerIDOrName": IDOrName,
			"message":           "容器删除失败",
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message":           "容器删除成功",
		"containerIDOrName": IDOrName,
	})
}

func containerSearch(c *gin.Context) {
	data := getJsonParam(c)
	IDOrName, err := data("IDOrName")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "json解析失败",
			"error":   err.Error(),
		})
		return
	}
	username, _ := c.Get("username")
	server := c.Query("server")
	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	containers, err := containers(true, username.(string), server)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取容器失败",
			"error":   err.Error(),
		})
		return
	}

	var searchedContainers []*Container
	for _, container := range containers {
		if container.Name == IDOrName || container.ID == IDOrName {
			searchedContainers = append(searchedContainers, container)
		}
	}
	if searchedContainers == nil {
		c.JSON(200, gin.H{
			"message": "未找到相关容器",
		})
		return
	}

	jsonData, err := json.Marshal(searchedContainers)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "将查找到的容器信息转换为json格式失败",
			"error":   err.Error(),
		})
		return
	}

	c.Data(200, "application/json", jsonData)
}
