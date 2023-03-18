package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const mockFileName = "a-mock-file"

func setLocalEnv(t *testing.T) {
	var err error
	// set mock env file
	err = os.Setenv("local_env", mockFileName)
	assert.NoError(t, err, "set env should no error")
}
func mockFileWithContent(t *testing.T, config string) {
	f, err := os.Create(mockFileName)
	assert.NoError(t, err, "create fail failed")

	_, err = f.WriteString(config)
	assert.NoError(t, err, "write file failed")
}

func removeMockFile(t *testing.T) {
	assert.NoError(t, os.Remove(mockFileName), "delete file failed, need to delete it manually")

}
