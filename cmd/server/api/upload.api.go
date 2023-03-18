package api

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/tony/mot-server/cmd/server/model"
	// "github.com/oklog/ulid/v2"
	"github.com/tony/mot-server/cmd/server/service"
)

type UploadAPI struct{}

type UploadParam struct {
	Filename string `json:"filename" binding:"required"`
	FilePath string `json:"filepath" binding:"required"`
}

func (a UploadAPI) UploadVideo(c *gin.Context) {
	log.Println("UploadVideo")
	f, err := c.FormFile("file")
	if err != nil {
		log.Printf("fail to read file from form: %s", err.Error())
		c.JSON(500, gin.H{"error": err.Error(), "message": "接收文件失败"})
		return
	}
	log.Println("uploading to minio")
	downloadUrl, err := UploadToMino(f)
	if err != nil {
		log.Printf("fail to upload to minio: %s", err.Error())
		c.JSON(500, gin.H{"error": err.Error(), "message": "接收文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully upload video", "url": downloadUrl.String()})
}

func UploadToMino(file *multipart.FileHeader) (*url.URL, error) {
	src, err := file.Open()
	if err != nil {
		log.Printf("fail to open file error: %s", err.Error())
		return nil, err
	}
	_, err = minioClient.PutObject("origin-video", file.Filename, src, -1, minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		log.Printf("fail to pub object error: %s", err.Error())
		return nil, err
	}

	reqParams := make(url.Values)
	expires := time.Second * 24 * 60 * 60
	presignedURL, err := minioClient.PresignedGetObject("origin-video", file.Filename, expires, reqParams)
	if err != nil {
		log.Printf("fail to pre signed url error: %s", err.Error())
		return nil, err
	}

	err = SaveToMysql(file.Filename, presignedURL.String())
	if err != nil {
		log.Printf("fail to save mysql error: %s", err.Error())
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
