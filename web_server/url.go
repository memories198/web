package web_server

import (
	"github.com/gin-gonic/gin"
	docker "web/docker_server"
	"web/user"
)

func registerUrl() {
	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/containers", containersList)
		authorized.GET("/containers/all", containersListAll)
		authorized.POST("/container/create", containerCreate)
		authorized.POST("/container/start", containerStart)
		authorized.POST("/container/stop", containerStop)
		authorized.POST("/container/kill", containerKill)
		authorized.POST("/container/remove", containerRemove)
		authorized.POST("/container/search", containerSearch)

		authorized.GET("/images", imageList)
		authorized.GET("/images/all", imageListAll)
		authorized.POST("/image/buildByFile", imageBuildByFile)
		authorized.POST("/image/tag/remove", imageTagRemove)
		authorized.POST("/image/removeByID", imageRemoveByID)
		authorized.POST("/image/tag/rename", imageTagRename)
		authorized.POST("/images/search/repository", imagesSearchRepository)
		authorized.POST("/images/search/local", imagesSearchLocal)
		authorized.POST("/image/addTag", imageAddTag)
		authorized.POST("/images/export", imagesExport)
		authorized.POST("/images/load", imagesLoad)
		authorized.POST("/image/pull", imagePull)
		authorized.POST("/image/push", imagePush)

		authorized.POST("/repository/login", repositoryLogin)
		authorized.POST("/repository/logout", repositoryLogout)
		authorized.GET("/repositories/list", repositoriesList)

		authorized.POST("/user/addServer", userAddServer)
		authorized.POST("/user/removeServer", userRemoveServer)
	}

	router.POST("/user/login", userLogin)
	router.POST("/user/register", userRegister)

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginTime, err := c.Cookie("username")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "用户未登录",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		username, err := user.MemoryGetKey(loginTime)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "用户未登录",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		err = setCookie(c, username)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "cookie设置失败",
			})
			c.Abort()
			return
		}

		clients, err := docker.InitUserClient(user.GetUserServers(username))
		if err != nil {
			c.JSON(400, gin.H{
				"message": "初始化用户docker客户端失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		userClients[username] = clients

		c.Set("username", username)

		c.Next()
	}
}
