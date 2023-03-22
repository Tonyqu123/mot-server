package service

import (
	// "github.com/gin-contrib/sessions"
	// "github.com/gin-gonic/gin"
	// "github.com/tony/mot-server/cmd/server/config"
	"github.com/tony/mot-server/cmd/server/model"
	// "github.com/tony/mot-server/cmd/server/utils"
)

// var rdb = config.InitRedis()

type LoginSrv struct {
	userName string
	password string
	role     string
}

type LoginBody struct {
	userName string `json:"username"`
	password string `json:"password"`
}

func (a LoginSrv) VerifyByUsername(username string, password string) (bool, *model.User) {
	userInfo, err := model.GetUserByUsername(username)

	if err != nil {
		return false, nil
	}

	// Check for username and password match, usually from a database
	if password != userInfo.Password {
		return false, nil
	}

	return true, userInfo
}
func (a LoginSrv) Register(user model.User) error {
	err := model.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}
