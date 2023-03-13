package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/tony/mot-server/cmd/server/api"
)

type Router struct {
	LoginAPI    api.LoginAPI
	FileAPI			api.FileAPI
	RabbitMQAPI	api.RabbitMQAPI
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
}