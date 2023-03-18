package config

type Minio struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyId     string `yaml:"access-key-id"`
	SecretAccessKey string `yaml:"secret-access-key"`
}

func GetMinioConfigOrDie() Minio {
	minioInfo := GetEnvOrDie().Db.Minio
	return minioInfo
}
