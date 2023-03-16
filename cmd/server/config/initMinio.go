package config

// func InitMinIO() {
// 	minioInfo := global.Settings.MinioInfo
// 	// 初使化 minio client对象。false是关闭https证书校验
// 	minioClient, err := minio.New(minioInfo.Endpoint, minioInfo.AccessKeyID, minioInfo.SecretAccessKey, false )
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//客户端注册到全局变量中
// 	global.MinioClient = minioClient
// 	//创建一个叫userheader的存储桶。
// 	// utils.CreateMinoBuket("userheader")
// }