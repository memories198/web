package web_server

import (
	"github.com/gin-gonic/gin"
	"web/dao"
	docker "web/docker_server"
)

func userAddServer(c *gin.Context) {
	data := getJsonParam(c)
	ipAndPort, err := data("ipAndPort")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "json解析失败",
			"error":   err.Error(),
		})
		return
	}
	username, _ := c.Get("username")
	err = dao.AddServer(username.(string), ipAndPort)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "保存docker服务器信息错误",
			"error":   err.Error(),
		})
		return
	}
	cli, err := docker.AddServerClient(ipAndPort)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "添加服务器至userClients失败",
			"error":   err.Error(),
		})
		return
	}
	userClients[username.(string)][ipAndPort] = cli

	c.JSON(200, gin.H{
		"message": "保存docker服务器信息成功",
	})
}
func userRemoveServer(c *gin.Context) {
	data := getJsonParam(c)
	ipAndPort, err := data("ipAndPort")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "json解析失败",
			"error":   err.Error(),
		})
		return
	}
	username, _ := c.Get("username")

	err = dao.RemoveServer(username.(string), ipAndPort)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "删除docker服务器信息失败",
			"error":   err.Error(),
		})
		return
	}
	delete(userClients[username.(string)], ipAndPort)
	c.JSON(200, gin.H{
		"message": "删除docker服务器信息成功",
	})
}
