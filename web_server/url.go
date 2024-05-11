package web_server

import (
	"github.com/gin-gonic/gin"
	"web/config"
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
		user.POST("/addServer", userAddServer)
		user.POST("/removeServer", userRemoveServer)
		user.GET("/listAllServer", userListAllServer)
	}
	router.POST("/user/login", userLogin)
	router.POST("/user/register", userRegister)
	router.NoRoute(func(c *gin.Context) {
		c.String(200, "hello")
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "用户未登录",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		jwt, err := VerifyJWT(token)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "无效的token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		username := jwt["username"].(string)

		id, err := dao.MemoryGetKey(username + "TokenUUID")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "获取" + username + "TokenUUID失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		} else if func() bool {
			if id == jwt["uuid"].(string) {
				return false
			}
			return true
		}() {
			if err != nil {
				c.JSON(400, gin.H{
					"message": "用户名或登录时间不正确",
				})
				c.Abort()
				return
			}
		}

		//更新cookie时间
		c.SetCookie("token", token, config.CookieExpireTime/1000, "/", "", false, true)
		err = dao.MemorySetExpire(username+"TokenUUID", 3600000)
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
		userClients = map[string]*docker.Client{}
		userClients[username.(string)] = cli
		c.Next()
	}
}
