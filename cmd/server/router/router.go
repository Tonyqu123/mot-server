package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tony/mot-server/cmd/server/api"
	"github.com/tony/mot-server/cmd/server/model"
	"net/http"
)

type Router struct {
	LoginAPI      api.LoginAPI
	FileAPI       api.FileAPI
	RabbitMQAPI   api.RabbitMQAPI
	UploadAPI     api.UploadAPI
	FileStatusAPI api.FileStatusAPI
}

func (a *Router) RegisterAPI(app *gin.Engine) {
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	app.GET("/get-file-list", a.FileAPI.GetFileList)

	app.POST("/login", a.LoginAPI.Login)

	app.GET("/send-mq", a.RabbitMQAPI.SendMessage)

	app.GET("/receive-mq", a.RabbitMQAPI.ReceiveMessage)

	app.GET("/get-minio", a.RabbitMQAPI.GetMinio)

	app.POST("/upload-video", a.UploadAPI.UploadVideo)

	app.GET("/test-product", model.GetProduct)

	app.POST("/insert-file-status", a.FileStatusAPI.AddFileStatus)
}
