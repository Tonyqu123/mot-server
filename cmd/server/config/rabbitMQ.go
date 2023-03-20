package config

type RabbitMQ struct {
	Endpoint string `yaml:"endpoint"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// GetRabbitMQFromEnv gets dsn string from the env
func GetRabbitMQFromEnv() RabbitMQ {
	rabbitMQInfo := GetEnvOrDie().Db.RabbitMQ
	//("amqp://admin:123456@localhost:5672/")
	return rabbitMQInfo
}
