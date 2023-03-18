package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockMinioConfig string = `
db:
  minio:
    endpoint: "127.0.0.1:9000"
    access-key-id: MOCK_ACCESS_KEY
    secret-access-key: CHANGEME123
`

func TestGetMinioConfigOrDie(t *testing.T) {
	// set mock env
	setLocalEnv(t)

	// create the mock file
	mockFileWithContent(t, mockMinioConfig)
	defer removeMockFile(t)

	cfg := GetMinioConfigOrDie()
	assert.Equal(t, cfg.AccessKeyId, "MOCK_ACCESS_KEY")
	assert.Equal(t, cfg.SecretAccessKey, "CHANGEME123")
	assert.Equal(t, cfg.Endpoint, "127.0.0.1:9000")
}
