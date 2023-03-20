package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockRabbitMQConfig string = `
db:
  rabbitMQ:
    endpoint: "localhost:5672"
    user: admin
    password: 123456
`

func TestGetRabbitMQFromEnv(t *testing.T) {
	// set mock env
	setLocalEnv(t)

	// create the mock file
	mockFileWithContent(t, mockRabbitMQConfig)
	defer removeMockFile(t)

	// execute the real func
	cfg := GetRabbitMQFromEnv()
	assert.Equal(t, cfg.Endpoint, "localhost:5672")
	assert.Equal(t, cfg.User, "admin")
	assert.Equal(t, cfg.Password, "123456")
}
