package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

//const configYaml = "/env/local.yaml"

type Env struct {
	Db struct {
		Mysql `yaml:"mysql"`
		Redis `yaml:"redis"`
		Minio `yaml:"minio"`
	} `yaml:"db"`
}

// if the envName specified, first read the value, use it if not empty;
// if empty then use the input Value.
// This is mainly to avoid the hard code.
func envFirst(envName, inputValue string) string {
	if os.Getenv(envName) != "" {
		return os.Getenv(envName)
	}
	return inputValue
}

func GetEnvOrDie() Env {
	t := Env{}
	data, err := os.ReadFile(envFirst("local_env", "/Users/litingting/GolandProjects/mot-server/cmd/server/env/local.yaml"))

	// 根据进程中的环境变量参数，判断使用哪份配置文件
	currentENV := os.Getenv("env")
	fmt.Println("currentENV", currentENV)
	if currentENV == "production" {
		data, err = os.ReadFile(GetAppPath() + "/env/prodEnv.yaml")
	} else if currentENV == "staging" {
		data, err = os.ReadFile(GetAppPath() + "/env/devEnv.yaml")
	}

	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		panic(err)
		return t
	}
	fmt.Println(t.Db.DSN())
	return t
}

// GetAppPath 解决执行 go build 时，获取不到相对路径的问题
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}
