package api

import (
	"fmt"
	"net/http"
	"time"

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

// 将上传成功的视频记录到数据库中
func SaveToMysql(filename string, originRrl string) error {
	var file model.File
	file.Filename = filename
	file.FileOrigin = originRrl
	// file.FileID = ulid.Make().String()
	file.FileTracked = ""
	file.UserID = 1
	file.UploadTime = time.Now().Format("2006-01-02 15:04:05")

	_, err := service.AddFile(file)
	if err != nil {
		return err
	}

	return nil
}
