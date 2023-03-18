package api

import (
	"fmt"
	"net/http"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tony/mot-server/cmd/server/model"
	"github.com/tony/mot-server/cmd/server/service"
	// "strings"
)

type FileAPI struct {
	FileSrv service.FileSrv
}

type FileParam struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a FileAPI) GetFileList(c *gin.Context) {
	var files []model.File
	files, count, err := a.FileSrv.GetFiles()
	fmt.Println("file api GetFiles count：", count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get files failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully list files", "data": files})
}
