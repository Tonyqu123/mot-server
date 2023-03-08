package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/tony/mot-server/cmd/server/api"
)

type Router struct {
	LoginAPI    api.LoginAPI
}

func (a *Router) RegisterAPI(app *gin.Engine) {
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})


	app.POST("/login", a.LoginAPI.Login)
}