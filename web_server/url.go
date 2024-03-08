package web_server

import (
	"github.com/gin-gonic/gin"
	"web/dao"
	docker "web/docker_server"
)

func registerUrl() {
	authorized := router.Group("/")
	authorized.Use(AuthMiddleware(), PingServer())
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

	}

	user := router.Group("/user")
	user.Use(AuthMiddleware())
	{
		user.POST("/user/addServer", userAddServer)
		user.POST("/user/removeServer", userRemoveServer)
		user.GET("/user/listAllServer", userListAllServer)
	}
	router.POST("/user/login", userLogin)
	router.POST("/user/register", userRegister)

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginTime, err := c.Cookie("loginTime")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "用户未登录",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		username, err := c.Cookie("username")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "用户未登录",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		lt, err := dao.MemoryGetKey(username + "LoginTime")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "获取用户登录时间失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		} else if lt != loginTime {
			if err != nil {
				c.JSON(400, gin.H{
					"message": "用户名或登录时间不正确",
				})
				c.Abort()
				return
			}
		}

		//更新cookie时间
		err = setCookie(c, username)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "更新cookie时间失败",
			})
			c.Abort()
			return
		}

		c.Set("username", username)

		c.Next()
	}
}
func PingServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		server := c.Query("server")
		cli, err := docker.PingServer(server)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "连接服务器失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		userClients = map[string]map[string]*docker.Client{username.(string): {}}
		userClients[username.(string)][server] = cli
		c.Next()
	}
}
