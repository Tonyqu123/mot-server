package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   int    `json:"status"` // 0 表示未通过申请，1 表示通过申请
	Files    []File `json:"files"`
}

func (User) TableName() string {
	return "user"
}

func GetUsers() ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return users, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := db.Where("username LIKE ?", username).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func UpdateUserStatusByUserid(id int) error {
	if err := db.Model(&User{}).Where("id = ?", id).Update("status", 1).Error; err != nil {
		return err
	}
	return nil
}

func DeleteByUserid(id int) error {
	// 根据主键删除
	if err := db.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByUsernameAndPass(username string) (User, error) {
	var user User
	if err := db.Where("name LIKE ?", username).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func CountUser() int64 {
	var total int64 = 0
	if err := db.Model(User{}).Count(&total).Error; err != nil {
		return -1
	}
	return total
}

func AddUser(user User) error {
	sql := db.Create(&user)
	if err := sql.Error; err != nil {
		return err
	}
	return nil
}
