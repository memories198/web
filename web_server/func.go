package web_server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	docker "web/docker_server"
	"web/user"
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

type Image struct {
	ID       string   `json:"ID"`
	Tags     []string `json:"Tags"`
	Created  string   `json:"Created"`
	Size     string   `json:"Size" `
	ParentID string   `json:"ParentID"`
}

func containersList(c *gin.Context) {
	username, _ := c.Get("username")
	server := c.Query("server")

	_, exist := userClients[username.(string)][server]
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

	_, exist := userClients[username.(string)][server]
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
	server := c.Query("server")

	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	ID, err := docker.CreateContainer(body, userClients[username.(string)][server])

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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.StartContainer(IDOrName, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.StopContainer(IDOrName, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.KillContainer(IDOrName, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.RemoveContainer(IDOrName, userClients[username.(string)][server])
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
	_, exist := userClients[username.(string)][server]
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

func imageList(c *gin.Context) {
	username, _ := c.Get("username")
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

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
	_, exist := userClients[username.(string)][server]
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	// 使用tar文件构建Docker镜像
	imageID, err := docker.BuildImage(saveFilePath, imageName, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.RemoveImage(imageName, userClients[username.(string)][server])
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
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	allImages, err := images(false, username.(string), server)
	for _, image := range allImages {
		if image.ID == imageID {
			err = docker.RemoveImage(imageID, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.RenameImage(imageTag, newTag, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	imagesInfo, err = docker.SearchImagesRepository(data["name"], userClients[username.(string)][server])
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
	_, exist := userClients[username.(string)][server]
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	imageID, err := docker.AddImageTag(imageName, newTag, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.ExportImages(c.Writer, images, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.LoadImages(saveFilePath, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.PullImage(image, userClients[username.(string)][server])
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
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.PushImage(image, userClients[username.(string)][server])
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

func repositoryLogin(c *gin.Context) {
	data := getJsonParam(c)
	repositoryUsername, err := data("username")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "json解析失败",
		})
		return
	}
	repositoryPassword, err := data("password")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "json解析失败",
		})
		return
	}
	serverAddress, err := data("serverAddress")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "json解析失败",
		})
		return
	}
	username, _ := c.Get("username")
	server := c.Query("server")
	_, exist := userClients[username.(string)][server]
	if exist == false {
		c.JSON(400, gin.H{
			"message": "docker服务器地址错误",
		})
		return
	}

	err = docker.LoginRepository(repositoryUsername, repositoryPassword, serverAddress, userClients[username.(string)][server])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "镜像仓库登录失败",
		})
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

	u := &user.User{
		Username: username,
		Password: password,
	}

	_, err = user.GetUser(username)
	if err == nil {
		c.JSON(400, gin.H{
			"message": "该用户已存在",
		})
		return
	}

	err = user.RegisterUser(u)
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

var userClients = map[string]map[string]*docker.Client{}
var cookieOutTime = 3600

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

	u, err := user.GetUser(username)
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
			})
			return
		}
	}

	clients, err := docker.InitUserClient(user.GetUserServers(username))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "初始化用户docker客户端失败",
			"error":   err.Error(),
		})
		return
	}
	userClients[username] = clients

	c.JSON(200, gin.H{
		"message": "登录成功",
	})
}

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
	err = user.AddServer(username.(string), ipAndPort)
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

	err = user.RemoveServer(username.(string), ipAndPort)
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
