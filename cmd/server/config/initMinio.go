package config

import (
	"fmt"
	"log"
	"github.com/minio/minio-go"
)

func InitMinio() *minio.Client{
	minioInfo := GetEnv().Db.Minio
	// 初使化 minio client对象。false是关闭https证书校验
	minioClient, err := minio.New(minioInfo.Endpoint, minioInfo.AccessKeyId, minioInfo.SecretAccessKey, false )
	if err != nil {
		log.Fatalln(err)
	}
	// 创建一个叫 originVideo 的存储桶。
	CreateMinoBuket(minioClient, "origin-video")
	CreateMinoBuket(minioClient, "tracked-video")
	return minioClient
}

// minio 初始化的时候 CreateMinoBuket 创建 两个 minio 桶，一个是源视频一个是跟踪的结果
func CreateMinoBuket(minioClient *minio.Client, bucketName string) {
	location := "us-east-1"
	err := minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := minioClient.BucketExists(bucketName)
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