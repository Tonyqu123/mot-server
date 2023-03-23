package api

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/tony/mot-server/cmd/server/config"
)

var (
	minioOnce   sync.Once
	minioClient *minio.Client
)

var VideoMap = map[string]string{
	"MOT16-02.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT16/tracked/MOT16-02.mp4",
	"MOT16-04.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT16/tracked/MOT16-04.mp4",
	"MOT16-05.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT16/tracked/MOT16-05.mp4",
	"MOT16-09.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT16/tracked/MOT16-09.mp4",
	"MOT16-10.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT16/tracked/MOT16-10.mp4",
	"MOT16-11.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT16/tracked/MOT16-11.mp4",
	"MOT16-13.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT16/tracked/MOT16-13.mp4",
	"MOT20-01.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT20/tracked/MOT20-01.mp4",
	"MOT20-02.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT20/tracked/MOT20-02.mp4",
	"MOT20-03.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT20/tracked/MOT20-03.mp4",
	"MOT20-05.mp4": "/Users/litingting/GolandProjects/mot-server/video/MOT20/tracked/MOT20-05.mp4",
}

type Form struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type Filename struct {
	Filename string `json:"filename" binding:"required"`
}

func InitMinioOrDie() error {
	minioOnce.Do(
		func() {
			minioClient = initMinioOrDie()

			// 创建一个叫 originVideo 的存储桶。
			CreateMinoBuket(minioClient, "origin-video")
			CreateMinoBuket(minioClient, "tracked-video")
		})
	return nil
}

func initMinioOrDie() *minio.Client {
	log.Println("initialing minIO")
	minioInfo := config.GetMinioConfigOrDie()
	// 初使化 minio client对象。false是关闭https证书校验
	mCli, err := minio.New(minioInfo.Endpoint, minioInfo.AccessKeyId, minioInfo.SecretAccessKey, false)
	if err != nil {
		log.Fatalf("FAIL TO intialize minIO, err: %s", err)
	}

	return mCli
}

type UploadAPI struct{}

type UploadParam struct {
	Filename string `json:"filename" binding:"required"`
	FilePath string `json:"filepath" binding:"required"`
}

func (a UploadAPI) UploadVideo(c *gin.Context) {
	//var form Form
	//err := c.ShouldBind(&form)
	//if err != nil {
	//	fmt.Println("err：", err)
	//}
	//
	//file, err := ioutil.ReadFile("/Users/litingting/Desktop/video/MOT20/origin/MOT20-03.mp4")
	//if err != nil {
	//	return
	//}
	//
	//downloadUrl, err := UploadToMino(f)
	//if err != nil {
	//	log.Printf("fail to upload to minio: %s", err.Error())
	//	c.JSON(500, gin.H{"error": err.Error(), "message": "接收文件失败"})
	//	return
	//}

	f, err := c.FormFile("file")

	fmt.Println("f.Filename：", f.Filename)

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

	str := fmt.Sprintf(f.Filename)
	fmt.Println("err = UploadToMinoByteType f.Filename：", str)
	fmt.Println("VideoMap[f.Filename]：", VideoMap[str])
	// 获取已经跟踪的本地文件，上传至 tracked-video
	err = UploadToMinoByteType(str)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "接收文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully upload video", "url": downloadUrl.String()})
}

func (a UploadAPI) DownloadTrackedUrl(c *gin.Context) {
	var filename Filename
	err := c.BindJSON(&filename)
	fmt.Println("filename：", filename.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不存在该文件！"})
		return
	}
	// 根据 bucket 名字和 文件名下载文件
	err = minioClient.FGetObject("tracked-video", filename.Filename, "/Users/litingting/Downloads/"+filename.Filename, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully download video"})
}

func UploadToMino(file *multipart.FileHeader) (*url.URL, error) {
	if minioClient == nil {
		log.Fatalln("minio client can't be nil")
	}

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

func UploadToMinoByteType(filename string) error {
	fmt.Println("UploadToMinoByteType filename：", filename)
	filePath := VideoMap[filename]
	fmt.Println("UploadToMinoByteType filePath：", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return err
	}

	uploadInfo, err := minioClient.PutObject("tracked-video", filename, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
	return nil
}

// minio 初始化的时候 CreateMinoBuket 创建 两个 minio 桶，一个是源视频一个是跟踪的结果
func CreateMinoBuket(cli *minio.Client, bucketName string) {
	location := "us-east-1"
	err := cli.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := cli.BucketExists(bucketName)
		fmt.Println(exists)
		if err == nil && exists {
			fmt.Printf("We already own %s\n", bucketName)
		} else {
			fmt.Println(err, exists)
			return
		}
	}

	fmt.Printf("Successfully created %s\n", bucketName)
}
