package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Router struct {

}

func (a *Router) RegisterAPI(app *gin.Engine) {
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}