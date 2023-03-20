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

type FileStatusAPI struct {
	FileStatusSrv service.FileStatusSrv
}

func (a FileStatusAPI) AddFileStatus(c *gin.Context) {
	var fileStatus model.FileStatus
	if err := c.Bind(&fileStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "参数有误"})
		return
	}

	fmt.Println("file_id：", fileStatus.FileID, "status：", fileStatus.Status)

	err := service.AddFileStatus(fileStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Insert file and status failed", "data": fileStatus.FileID})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully insert file and status", "data": fileStatus.FileID})
}
