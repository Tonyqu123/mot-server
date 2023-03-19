package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockMySqlConfig string = `
db:
  mysql:
    host: "my.host"
    port: 3306
    user: "root"
    password: "123456"
    db-name: "dbName"
    parameters: "foo=bar"
`

func TestGetDsnFromEnv(t *testing.T) {
	// set mock env
	setLocalEnv(t)

	// create the mock file
	mockFileWithContent(t, mockMySqlConfig)
	defer removeMockFile(t)

	// execute the real func
	dsn := GetDsnFromEnv()
	assert.Equal(t, dsn, "root:123456@tcp(my.host:3306)/dbName?foo=bar")

}
