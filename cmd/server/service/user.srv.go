package service

import (
	// "encoding/json"
	"github.com/tony/mot-server/cmd/server/model"
)

type UserSrv struct{}

func GetUsers() ([]model.User, error) {
	users, err := model.GetUsers()
	//total := model.CountUser()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func AddUser(user model.User) error {
	err := model.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

func PermitUser(id int) error {
	err := model.UpdateUserStatusByUserid(id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(id int) error {
	err := model.DeleteByUserid(id)
	if err != nil {
		return err
	}
	return nil
}
