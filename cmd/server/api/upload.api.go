package api

import (
	"fmt"
	"time"
	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"net/http"
	"net/url"
	"mime/multipart"
	// "github.com/oklog/ulid/v2"
	"github.com/tony/mot-server/cmd/server/service"
	"github.com/tony/mot-server/cmd/server/model"
)

type UploadAPI struct {}

type UploadParam struct {
	Filename string `json:"filename" binding:"required"`
	FilePath string `json:"filepath" binding:"required"`
}


func (a UploadAPI) UploadVideo(c *gin.Context) {
	f, err := c.FormFile("file")
	fmt.Println("UploadVideo")
	if err != nil {
		 c.String(http.StatusBadRequest, "接收文件失败")
	}
	downloadUrl, err := UploadToMino(f)

	if err != nil {
		c.String(http.StatusBadRequest, "接收文件失败")
 	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully upload video", "url": downloadUrl.String()})
}

func UploadToMino(file *multipart.FileHeader) (*url.URL, error){
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	_, err = minioClient.PutObject("origin-video", file.Filename, src,-1, minio.PutObjectOptions{ContentType:"video/mp4"})
	if err != nil {
		return nil, err
	}

	reqParams := make(url.Values)
	expires := time.Second*24*60*60
	presignedURL, err := minioClient.PresignedGetObject("origin-video", file.Filename, expires, reqParams)
	if err != nil {
		return nil, err
	}

	err = SaveToMysql(file.Filename, presignedURL.String())

	if err != nil {
		return presignedURL, err
	}

	fmt.Println("presignedURL：", presignedURL)

	return presignedURL, nil
}

// 将上传成功的视频记录到数据库中
func SaveToMysql(filename string, originRrl string) error {
	var file model.File
	file.Filename = filename
	file.FileOrigin = originRrl
	// file.Fileid = ulid.Make().String()
	file.FileTracked = ""
	file.Userid = "1"
	file.Uploadtime = time.Now().Format("2006-01-02 15:04:05")

	_, err := service.AddFile(file)
	if err != nil {
		return err
	}

	return nil
}