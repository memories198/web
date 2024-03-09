package web_server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"web/dao"
	docker "web/docker_server"
)

var userAuthConfigurations = map[string][]docker.AuthConfiguration{}

func repositoriesList(c *gin.Context) {
	repositories := docker.ListRepositories()
	if repositories == nil {
		c.JSON(200, gin.H{
			"message": "没有登录任何镜像仓库",
		})
		return
	}
	jsonData, err := json.Marshal(repositories)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "将获取到的仓库信息转换为json格式错误",
			"error":   err.Error(),
		})
		return
	}
	c.Data(200, "application/json", jsonData)
}

func repositoryLogin(c *gin.Context) {
	data := getJsonParam(c)
	repositoryUsername, err := data("repositoryUsername")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "json解析失败",
		})
		return
	}
	repositoryPassword, err := data("repositoryPassword")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "json解析失败",
		})
		return
	}
	repository, err := data("repository")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "json解析失败",
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

	err = docker.LoginRepository(repositoryUsername, repositoryPassword, repository, userClients[username.(string)])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "镜像仓库登录失败",
		})
		return
	}

	err = dao.AddRepository(username.(string), repositoryUsername, repositoryPassword)
	if err != nil {

		return
	}
	c.JSON(200, gin.H{
		"message": "镜像仓库登录成功",
	})
}
func repositoryLogout(c *gin.Context) {
	data := getJsonParam(c)
	repository, err := data("repository")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "json解析失败",
		})
		return
	}

	err = docker.LogoutRepository(repository)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "注销镜像仓库失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "注销镜像仓库成功",
	})
}
