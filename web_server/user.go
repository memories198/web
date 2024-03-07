package web_server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"web/dao"
	docker "web/docker_server"
)

var userClients = map[string]map[string]*docker.Client{}

func userLogin(c *gin.Context) {
	data := getJsonParam(c)
	username, err := data("username")
	if err != nil {
		return
	}
	password, err := data("password")
	if err != nil {
		return
	}

	u, err := dao.GetUser(username)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "从数据库中读取用户信息失败",
			"error":   err.Error(),
		})
		return
	} else if u.Username == "" || u.Password != password {
		c.JSON(400, gin.H{
			"message": "登录失败，用户名或密码不正确",
		})
		return
	} else {
		err = setCookie(c, username)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "cookie设置失败",
				"error":   err.Error(),
			})
			return
		}
	}

	var errs []string
	userClients[username], errs = docker.InitUserClient(dao.GetUserServers(username))

	marshal, err := json.Marshal(struct {
		Err []string `json:"errors"`
	}{errs})
	if err != nil {
		c.JSON(400, gin.H{
			"message": "将errs转换为json格式失败",
			"error":   err.Error(),
		})
		return
	}
	c.Data(200, "application/json", marshal)
}

func userRegister(c *gin.Context) {
	data := getJsonParam(c)
	username, err := data("username")
	if err != nil {
		return
	}
	password, err := data("password")
	if err != nil {
		return
	}

	u := &dao.User{
		Username: username,
		Password: password,
	}

	_, err = dao.GetUser(username)
	if err == nil {
		c.JSON(400, gin.H{
			"message": "该用户已存在",
		})
		return
	}

	err = dao.RegisterUser(u)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "注册失败，保存至数据库失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "注册成功",
	})
}
