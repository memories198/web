package web_server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	docker "web/docker_server"
)

type Image struct {
	ID       string   `json:"ID"`
	Tags     []string `json:"Tags"`
	Created  string   `json:"Created"`
	Size     string   `json:"Size" `
	ParentID string   `json:"ParentID"`
}

func imageList(c *gin.Context) {
	username, _ := c.Get("username")
	server := c.Query("server")

	imagesInfo, err := images(false, username.(string), server)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "镜像获取失败",
			"error":   err.Error(),
		})
		return
	}
	marshal, err := json.Marshal(imagesInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "查看镜像失败",
			"error":   err.Error(),
		})
	}
	c.Data(200, "application/json", marshal)
}
func imageListAll(c *gin.Context) {
	username, _ := c.Get("username")
	server := c.Query("server")
	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	imagesInfo, err := images(true, username.(string), server)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "镜像获取失败",
			"error":   err.Error(),
		})
		return
	}
	marshal, err := json.Marshal(imagesInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "将查找到的镜像转换为json失败",
			"error":   err.Error(),
		})
	}
	c.Data(200, "application/json", marshal)
}

func imageBuildByFile(c *gin.Context) {
	imageName := c.PostForm("imageName")
	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":   "imagesName格式错误",
			"imageName": imageName,
		})
		return
	}

	file, err := c.FormFile("Dockerfile")
	if file == nil {
		c.JSON(400, gin.H{
			"message": "未上传文件",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"file":    file.Filename,
			"message": "文件解析出错",
		})
		return
	}

	if len(file.Filename) < 5 || file.Filename[len(file.Filename)-4:] != ".tar" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文件格式错误",
			"file":    file.Filename,
		})
		return
	}
	unixNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	saveFilePath := "./web_server/upload/" + unixNano + file.Filename
	err = c.SaveUploadedFile(file, saveFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "文件保存出错",
			"file":    file.Filename,
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

	// 使用tar文件构建Docker镜像
	imageID, err := docker.BuildImage(saveFilePath, imageName, userClients[username.(string)])
	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "镜像构建失败",
			"image":   imageName,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "镜像构建成功",
		"image":   imageName,
		"imageID": imageID,
	})
}
func imageTagRemove(c *gin.Context) {
	data := getJsonParam(c)
	imageName, err := data("imageName")
	if err != nil {
		c.JSON(400, gin.H{
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

	err = docker.RemoveImage(imageName, userClients[username.(string)])
	if err != nil {
		c.JSON(400, gin.H{
			"message": "镜像删除失败",
			"error":   err.Error(),
			"image":   imageName,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "镜像删除成功",
		"image":   imageName,
	})
}
func imageRemoveByID(c *gin.Context) {
	data := getJsonParam(c)
	imageID, err := data("ID")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "json解析失败",
			"err":     err.Error(),
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

	allImages, err := images(false, username.(string), server)
	for _, image := range allImages {
		if image.ID == imageID {
			err = docker.RemoveImage(imageID, userClients[username.(string)])
			if err != nil {
				c.JSON(400, gin.H{
					"message": "镜像删除失败",
					"err":     err.Error(),
				})
				return
			}
		}
	}

	c.JSON(200, gin.H{
		"message": "镜像删除成功",
		"imageID": imageID,
	})
}
func imageTagRename(c *gin.Context) {
	data := getJsonParam(c)
	imageTag, err := data("imageIDOrName")
	if imageTag == "" || err != nil {
		c.JSON(400, gin.H{
			"message": "imageTag不合法",
			"err":     err.Error(),
		})
		return
	}
	newTag, err := data("newTag")
	if newTag == "" || err != nil {
		c.JSON(400, gin.H{
			"message": "newTag不合法",
			"err":     err.Error(),
		})
		return
	}
	if imageTag == newTag {
		c.JSON(400, gin.H{
			"message": "镜像重命名出错",
			"error":   "旧的的名称与新的名称相同",
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

	err = docker.RenameImage(imageTag, newTag, userClients[username.(string)])
	if err != nil {
		c.JSON(400, gin.H{
			"message": "镜像重命名出错",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "镜像重命名成功",
		"newTag":  newTag,
	})
}

func imagesSearchRepository(c *gin.Context) {
	var data map[string]string

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "json解析出错",
			"error":   err.Error(),
		})
		return
	}

	var imagesInfo []*docker.ImageInfo
	username, _ := c.Get("username")
	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	imagesInfo, err = docker.SearchImagesRepository(data["name"], userClients[username.(string)])
	if err != nil {
		c.JSON(400, gin.H{
			"message": "查找镜像失败",
			"error":   err.Error(),
		})
		return
	} else if imagesInfo == nil {
		c.JSON(200, gin.H{
			"message": "未找到相关镜像",
		})
		return
	}

	imagesInfoJson, err := json.Marshal(imagesInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "将查找到的镜像信息转换为json失败",
			"error":   err.Error(),
		})
		return
	}
	c.Data(200, "application/json", imagesInfoJson)
}
func imagesSearchLocal(c *gin.Context) {
	data := getJsonParam(c)
	name, err := data("imageIDOrName")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "解析json失败",
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

	imagesInfo, err := images(true, username.(string), server)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取所有镜像信息失败",
			"error":   err.Error(),
		})
		return
	}

	var images []Image
loop1:
	for _, image := range imagesInfo {
		if image.ID == name {
			images = append(images, image)
			break loop1
		}
		for _, imageName := range image.Tags {
			n := strings.Split(imageName, ":")
			if imageName == name { //镜像的名称和ID是唯一的，所以只需找到一个就要跳出循环
				images = append(images, image)
				break loop1
			}
			if n[0] == name {
				images = append(images, image)
			}
		}
	}

	if images == nil {
		c.JSON(200, gin.H{
			"message": "未查找到相关镜像",
		})
		return
	}

	ans, err := json.Marshal(images)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "将查找到的镜像信息转换为json失败",
			"error":   err.Error(),
		})
		return
	}
	c.Data(200, "application/json", ans)
}
func imageAddTag(c *gin.Context) {
	data := getJsonParam(c)
	imageName, err := data("imageIDOrName")
	newTag, err := data("tag")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "解析json失败",
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

	imageID, err := docker.AddImageTag(imageName, newTag, userClients[username.(string)])
	if err != nil {
		c.JSON(400, gin.H{
			"message": "容器标签添加失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message":     "容器标签添加成功",
		"imageID":     imageID,
		"imageNewTag": newTag,
	})
}

func imagesExport(c *gin.Context) {
	data := getJsonAnyParam(c)
	fileNameEmptyInterface, _ := data("fileName")
	var fileName string
	if fileNameEmptyInterface == nil {
		fileName = "images.tar"
	} else {
		switch fileNameEmptyInterface.(type) {
		case string:
			if fileNameEmptyInterface.(string) == "" {
				fileName = "images.tar"
			} else {
				fileName = fileNameEmptyInterface.(string)
			}
		}
	}

	imagesEmptyInterface, err := data("images")
	var images []string
	if err != nil {
		c.JSON(400, gin.H{
			"message": "json解析失败",
			"error":   err.Error(),
		})
		return
	} else if imagesEmptyInterface == nil {
		c.JSON(400, gin.H{
			"message": "json解析失败",
			"error":   "缺少镜像名称",
		})
		return
	} else {
		for _, i := range imagesEmptyInterface.([]interface{}) {
			switch i.(type) {
			case string:
				images = append(images, i.(string))
			default:
				c.JSON(400, gin.H{
					"message": "json解析失败",
					"error":   "images字段格式错误",
				})
				return
			}
		}
	}

	// 设置响应头信息以下载文件，而不是在浏览器中打开
	// 这里的filename可以根据实际情况修改
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")

	username, _ := c.Get("username")
	_, exist := userClients[username.(string)]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.ExportImages(c.Writer, images, userClients[username.(string)])
	if err != nil {
		log.Println(err)
	}
}

func imagesLoad(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "json解析失败",
			"error":   err.Error(),
		})
		return
	}

	if len(file.Filename) < 5 || file.Filename[len(file.Filename)-4:] != ".tar" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文件格式错误,文件格式需要为.tar",
			"file":    file.Filename,
		})
		return
	}
	unixNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	saveFilePath := "./web_server/upload/" + unixNano + file.Filename
	err = c.SaveUploadedFile(file, saveFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "文件保存出错",
			"file":    file.Filename,
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

	err = docker.LoadImages(saveFilePath, userClients[username.(string)])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "镜像加载出错",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "镜像加载成功",
	})
}
func imagePull(c *gin.Context) {
	data := getJsonParam(c)
	image, err := data("image")
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

	err = docker.PullImage(image, userClients[username.(string)])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "拉取镜像失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "拉取镜像成功",
	})
}

func imagePush(c *gin.Context) {
	data := getJsonParam(c)
	image, err := data("image")
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

	err = docker.PushImage(image, userClients[username.(string)])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "镜像推送失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "镜像推送成功",
	})
}
