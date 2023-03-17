package api

import (
	"fmt"
	// "github.com/gin-contrib/sessions"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"net/http"
	"net/url"
	"mime/multipart"
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

	fmt.Println("presignedURL：", presignedURL)

	return presignedURL, nil
}